package sentinelone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

const (
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 3
	apiPathPrefix     = "/web/api/v2.1"
	authHeaderName    = "Authorization"
	authHeaderPrefix  = "ApiToken "
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
	errTag            = "sentinelone: "
)

// Client is a SentinelOne Management API client.
//
// API calls are grouped by resource type and accessed through sub-client
// fields: [Client.Accounts], [Client.Sites], [Client.RBAC], and
// [Client.Users].  All methods accept a [context.Context] as their first
// argument, which controls request cancellation and deadline propagation.
// Non-2xx responses are returned as *[types.ResponseError].
// Create an instance with [NewClient].
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
	rateLimits []pathRateLimit // nil = disabled; sorted longest-prefix-first
	maxRetries int             // max 429 retry attempts; default 3

	Accounts *AccountsClient
	Sites    *SitesClient
	RBAC     *RBACClient
	Users    *UsersClient
	Agents   *AgentsClient
	Licenses *LicensesClient
}

// ClientOption is a functional option passed to [NewClient] to customise
// the underlying HTTP client behaviour.
type ClientOption func(*Client)

func (ClientOption) applyLoadOption() {}

// WithHTTPClient replaces the default [net/http.Client] used for requests.
// This is the right hook for custom transports, proxies, or mTLS configuration.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) { c.httpClient = hc }
}

// WithTimeout overrides the default 30-second per-request timeout.
// This sets the Timeout field on the underlying [net/http.Client]; it is
// superseded by any deadline already present on the request [context.Context].
func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// WithRateLimiting enables or disables the built-in per-path token-bucket
// rate limiter.  It is enabled by default using the official SentinelOne rate
// limits.  Disable it when you manage throttling externally or wrap the client
// behind your own middleware.
func WithRateLimiting(enabled bool) ClientOption {
	return func(c *Client) {
		if enabled {
			c.rateLimits = buildDefaultRateLimits()
		} else {
			c.rateLimits = nil
		}
	}
}

// WithMaxRetries sets the maximum number of times the client will retry a
// request that received a 429 Too Many Requests response.  The default is 3.
// Set to 0 to disable retries entirely.
func WithMaxRetries(n int) ClientOption {
	return func(c *Client) { c.maxRetries = n }
}

// NewClient creates a new SentinelOne Management API client.
//
// baseURL is the root URL of the management console, for example
// "https://your-tenant.sentinelone.net".  A trailing slash is stripped
// automatically.
//
// apiToken is sent as "ApiToken <token>" in every request's Authorization
// header.  Tokens can be generated in the console (My User → API Token
// Operations) or via [UsersClient.GenerateAPIToken].
//
// Additional behaviour can be tuned with [WithTimeout] and [WithHTTPClient].
func NewClient(baseURL, apiToken string, opts ...ClientOption) *Client {
	cli := &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		rateLimits: buildDefaultRateLimits(),
		maxRetries: defaultMaxRetries,
	}
	for _, o := range opts {
		o(cli)
	}

	cli.Accounts = &AccountsClient{c: cli}
	cli.Sites = &SitesClient{c: cli}
	cli.RBAC = &RBACClient{c: cli}
	cli.Users = &UsersClient{c: cli}
	cli.Agents = &AgentsClient{c: cli}
	cli.Licenses = &LicensesClient{c: cli}

	return cli
}

// AsResponseError unwraps err into a *[types.ResponseError] and reports
// whether the conversion succeeded.  Use this instead of errors.As when you
// want both the typed value and the ok-check in one call:
//
//	if respErr, ok := sentinelone.AsResponseError(err); ok {
//	    fmt.Println(respErr.StatusCode, respErr.Errors)
//	}
func AsResponseError(err error) (*types.ResponseError, bool) {
	return types.AsResponseError(err)
}

// rawResponse is the generic envelope for all API responses.
type rawResponse struct {
	Pagination *types.Pagination `json:"pagination,omitempty"`
	Errors     []json.RawMessage `json:"errors,omitempty"`
	Data       json.RawMessage   `json:"data,omitempty"`
}

// buildURL constructs the full request URL with optional query parameters.
func (c *Client) buildURL(path string, params url.Values) string {
	u := c.baseURL + apiPathPrefix + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	return u
}

// do executes an HTTP request and decodes the response envelope.
// It applies the proactive per-path rate limiter before each attempt and
// retries automatically on 429 responses (up to c.maxRetries times), honouring
// the Retry-After header when present.
func (c *Client) do(ctx context.Context, method, path string, params url.Values, body, out any) ( //nolint:funlen
	*types.Pagination, error,
) {
	// Marshal the body once; each retry creates a fresh bytes.Reader from the
	// same slice so the reader position is always at the start.
	var bodyBytes []byte

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("%s marshal request: %w", errTag, err)
		}

		bodyBytes = b
	}

	for attempt := 0; ; attempt++ {
		// Proactive rate limiting — block until the limiter grants a token.
		lim := c.limiterFor(path)

		if lim != nil {
			err := lim.Wait(ctx)
			if err != nil {
				return nil, fmt.Errorf("%s rate limit wait: %w", errTag, err)
			}
		}

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, c.buildURL(path, params), bodyReader)
		if err != nil {
			return nil, fmt.Errorf("%s build request: %w", errTag, err)
		}

		req.Header.Set(authHeaderName, authHeaderPrefix+c.apiToken)

		if bodyBytes != nil {
			req.Header.Set(contentTypeHeader, contentTypeJSON)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("%s HTTP request: %w", errTag, err)
		}

		// Read and close the body explicitly — defer inside a loop would
		// accumulate closers until the function returns.
		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if readErr != nil {
			return nil, fmt.Errorf("%s read response: %w", errTag, readErr)
		}

		// Reactive rate limiting — honour 429 and retry up to maxRetries times.
		if resp.StatusCode == http.StatusTooManyRequests && attempt < c.maxRetries {
			wait := retryAfterDuration(resp)

			select {
			case <-time.After(wait):
				continue
			case <-ctx.Done():
				return nil, fmt.Errorf("%s context done during 429 backoff: %w", errTag, ctx.Err())
			}
		}

		raw, err := c.decodeRawResponse(respBody, out)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, &types.ResponseError{StatusCode: resp.StatusCode, Errors: parseErrors(raw.Errors)}
		}

		return raw.Pagination, nil
	}
}

// decodeRawResponse unmarshals the API envelope and, when out is non-nil,
// further unmarshals the data field into out.
func (c *Client) decodeRawResponse(body []byte, out any) (rawResponse, error) {
	var raw rawResponse

	if len(body) > 0 {
		err := json.Unmarshal(body, &raw)
		if err != nil {
			return raw, fmt.Errorf("%s decode response: %w", errTag, err)
		}
	}

	if out != nil && len(raw.Data) > 0 {
		err := json.Unmarshal(raw.Data, out)
		if err != nil {
			return raw, fmt.Errorf("%s decode data: %w", errTag, err)
		}
	}

	return raw, nil
}

func parseErrors(raw []json.RawMessage) []types.APIError {
	if len(raw) == 0 {
		return nil
	}

	errs := make([]types.APIError, 0, len(raw))
	for _, r := range raw {
		var e types.APIError
		if json.Unmarshal(r, &e) == nil {
			errs = append(errs, e)
		}
	}

	return errs
}

// get performs a GET request.
func (c *Client) get(ctx context.Context, path string, params url.Values, out any) (*types.Pagination, error) {
	return c.do(ctx, http.MethodGet, path, params, nil, out)
}

// post performs a POST request.
func (c *Client) post(ctx context.Context, path string, body, out any) (*types.Pagination, error) {
	return c.do(ctx, http.MethodPost, path, nil, body, out)
}

// put performs a PUT request.
func (c *Client) put(ctx context.Context, path string, body, out any) (*types.Pagination, error) {
	return c.do(ctx, http.MethodPut, path, nil, body, out)
}

// delete performs a DELETE request.
func (c *Client) delete(ctx context.Context, path string, body any) (*types.Pagination, error) {
	return c.do(ctx, http.MethodDelete, path, nil, body, nil)
}

// -- Query param helpers --

// ListParams contains the common pagination and sorting parameters shared by
// every list endpoint.  Embed this struct in the endpoint-specific params
// type (e.g. [ListAccountsParams]) — all fields are optional.
//
//   - Skip: number of records to skip (0–1000).  Use Cursor for deeper pages.
//   - Limit: maximum records to return per page (1–1000; API default is 10).
//   - Cursor: opaque cursor returned in the previous [types.Pagination].NextCursor.
//   - CountOnly: when true, the response contains only TotalItems — no data.
//   - SkipCount: when true, TotalItems is not calculated (faster for large sets).
//   - SortBy: field name to sort by (valid values differ per endpoint).
//   - SortOrder: "asc" or "desc" (default is "asc" if SortBy is set).
type ListParams struct {
	Skip      *int
	Limit     *int
	Cursor    *string
	CountOnly *bool
	SkipCount *bool
	SortBy    *string
	SortOrder *string
}

func (p *ListParams) values() url.Values {
	vals := url.Values{}
	if p == nil {
		return vals
	}

	if p.Skip != nil {
		vals.Set("skip", strconv.Itoa(*p.Skip))
	}

	if p.Limit != nil {
		vals.Set("limit", strconv.Itoa(*p.Limit))
	}

	if p.Cursor != nil {
		vals.Set("cursor", *p.Cursor)
	}

	if p.CountOnly != nil {
		vals.Set("countOnly", strconv.FormatBool(*p.CountOnly))
	}

	if p.SkipCount != nil {
		vals.Set("skipCount", strconv.FormatBool(*p.SkipCount))
	}

	if p.SortBy != nil {
		vals.Set("sortBy", *p.SortBy)
	}

	if p.SortOrder != nil {
		vals.Set("sortOrder", *p.SortOrder)
	}

	return vals
}

func setBool(v url.Values, key string, val *bool) {
	if val != nil {
		v.Set(key, strconv.FormatBool(*val))
	}
}

func setString(v url.Values, key string, val *string) {
	if val != nil && *val != "" {
		v.Set(key, *val)
	}
}

func setStringSlice(v url.Values, key string, vals []string) {
	if len(vals) > 0 {
		v.Set(key, strings.Join(vals, ","))
	}
}

func setInt(v url.Values, key string, val *int) {
	if val != nil {
		v.Set(key, strconv.Itoa(*val))
	}
}

// BoolPtr returns a pointer to the given bool value.
// Useful for optional *bool fields on request structs.
func BoolPtr(v bool) *bool { return &v }

// IntPtr returns a pointer to the given int value.
// Useful for optional *int fields such as ListParams.Limit.
func IntPtr(v int) *int { return &v }

// StringPtr returns a pointer to the given string value.
// Useful for optional *string fields on request and filter structs.
func StringPtr(v string) *string { return &v }
