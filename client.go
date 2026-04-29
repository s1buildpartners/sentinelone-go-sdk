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
)

const (
	defaultTimeout    = 30 * time.Second
	apiPathPrefix     = "/web/api/v2.1"
	authHeaderName    = "Authorization"
	authHeaderPrefix  = "ApiToken "
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

// Client is a SentinelOne Management API client.
//
// All methods accept a [context.Context] as their first argument, which
// controls request cancellation and deadline propagation.  Non-2xx responses
// are returned as *[ResponseError].  Create an instance with [NewClient].
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

// ClientOption is a functional option passed to [NewClient] to customise
// the underlying HTTP client behaviour.
type ClientOption func(*Client)

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

// NewClient creates a new SentinelOne Management API client.
//
// baseURL is the root URL of the management console, for example
// "https://your-tenant.sentinelone.net".  A trailing slash is stripped
// automatically.
//
// apiToken is sent as "ApiToken <token>" in every request's Authorization
// header.  Tokens can be generated in the console (My User → API Token
// Operations) or via [Client.GenerateAPIToken].
//
// Additional behaviour can be tuned with [WithTimeout] and [WithHTTPClient].
func NewClient(baseURL, apiToken string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Pagination holds cursor-based pagination metadata returned by list endpoints.
//
// When NextCursor is non-nil, pass its value as ListParams.Cursor on the next
// call to retrieve the following page.  NextCursor is nil on the last page.
// TotalItems reports the total number of matching records across all pages
// (unless the request set SkipCount, in which case it may be zero).
type Pagination struct {
	NextCursor *string `json:"nextCursor,omitempty"`
	TotalItems int     `json:"totalItems"`
}

// APIError represents a single structured error object returned by the API
// inside the top-level "errors" array of the response envelope.
type APIError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func (e APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("sentinelone: %s (%s)", e.Message, e.Detail)
	}
	return fmt.Sprintf("sentinelone: %s", e.Message)
}

// ResponseError is returned when the server responds with a non-2xx status
// code.  Use [errors.As] to unwrap it and inspect StatusCode or Errors:
//
//	var respErr *sentinelone.ResponseError
//	if errors.As(err, &respErr) {
//	    fmt.Println(respErr.StatusCode, respErr.Errors)
//	}
type ResponseError struct {
	StatusCode int
	Errors     []APIError
}

func (e *ResponseError) Error() string {
	if len(e.Errors) > 0 {
		msgs := make([]string, len(e.Errors))
		for i, err := range e.Errors {
			msgs[i] = err.Error()
		}
		return fmt.Sprintf("sentinelone: HTTP %d: %s", e.StatusCode, strings.Join(msgs, "; "))
	}
	return fmt.Sprintf("sentinelone: HTTP %d", e.StatusCode)
}

// rawResponse is the generic envelope for all API responses.
type rawResponse struct {
	Pagination *Pagination       `json:"pagination,omitempty"`
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
func (c *Client) do(ctx context.Context, method, path string, params url.Values, body interface{}, out interface{}) (*Pagination, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("sentinelone: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(path, params), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("sentinelone: build request: %w", err)
	}
	req.Header.Set(authHeaderName, authHeaderPrefix+c.apiToken)
	if body != nil {
		req.Header.Set(contentTypeHeader, contentTypeJSON)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sentinelone: HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sentinelone: read response: %w", err)
	}

	var raw rawResponse
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &raw); err != nil {
			return nil, fmt.Errorf("sentinelone: decode response: %w", err)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErrs := parseErrors(raw.Errors)
		return nil, &ResponseError{StatusCode: resp.StatusCode, Errors: apiErrs}
	}

	if out != nil && len(raw.Data) > 0 {
		if err := json.Unmarshal(raw.Data, out); err != nil {
			return nil, fmt.Errorf("sentinelone: decode data: %w", err)
		}
	}

	return raw.Pagination, nil
}

func parseErrors(raw []json.RawMessage) []APIError {
	if len(raw) == 0 {
		return nil
	}
	errs := make([]APIError, 0, len(raw))
	for _, r := range raw {
		var e APIError
		if json.Unmarshal(r, &e) == nil {
			errs = append(errs, e)
		}
	}
	return errs
}

// get performs a GET request.
func (c *Client) get(ctx context.Context, path string, params url.Values, out interface{}) (*Pagination, error) {
	return c.do(ctx, http.MethodGet, path, params, nil, out)
}

// post performs a POST request.
func (c *Client) post(ctx context.Context, path string, body interface{}, out interface{}) (*Pagination, error) {
	return c.do(ctx, http.MethodPost, path, nil, body, out)
}

// put performs a PUT request.
func (c *Client) put(ctx context.Context, path string, body interface{}, out interface{}) (*Pagination, error) {
	return c.do(ctx, http.MethodPut, path, nil, body, out)
}

// delete performs a DELETE request.
func (c *Client) delete(ctx context.Context, path string, body interface{}, out interface{}) (*Pagination, error) {
	return c.do(ctx, http.MethodDelete, path, nil, body, out)
}

// -- Query param helpers --

// ListParams contains the common pagination and sorting parameters shared by
// every list endpoint.  Embed this struct in the endpoint-specific params
// type (e.g. [ListAccountsParams]) — all fields are optional.
//
//   - Skip: number of records to skip (0–1000).  Use Cursor for deeper pages.
//   - Limit: maximum records to return per page (1–1000; API default is 10).
//   - Cursor: opaque cursor returned in the previous [Pagination].NextCursor.
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
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.Skip != nil {
		v.Set("skip", strconv.Itoa(*p.Skip))
	}
	if p.Limit != nil {
		v.Set("limit", strconv.Itoa(*p.Limit))
	}
	if p.Cursor != nil {
		v.Set("cursor", *p.Cursor)
	}
	if p.CountOnly != nil {
		v.Set("countOnly", strconv.FormatBool(*p.CountOnly))
	}
	if p.SkipCount != nil {
		v.Set("skipCount", strconv.FormatBool(*p.SkipCount))
	}
	if p.SortBy != nil {
		v.Set("sortBy", *p.SortBy)
	}
	if p.SortOrder != nil {
		v.Set("sortOrder", *p.SortOrder)
	}
	return v
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

func setInt(v url.Values, key string, val *int) {
	if val != nil {
		v.Set(key, strconv.Itoa(*val))
	}
}

func setStringSlice(v url.Values, key string, vals []string) {
	if len(vals) > 0 {
		v.Set(key, strings.Join(vals, ","))
	}
}

func ptr[T any](v T) *T { return &v }

// BoolPtr returns a pointer to the given bool value.
// Useful for optional *bool fields on request structs.
func BoolPtr(v bool) *bool { return &v }

// IntPtr returns a pointer to the given int value.
// Useful for optional *int fields such as ListParams.Limit.
func IntPtr(v int) *int { return &v }

// StringPtr returns a pointer to the given string value.
// Useful for optional *string fields on request and filter structs.
func StringPtr(v string) *string { return &v }
