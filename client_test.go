package sentinelone

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// ---- helpers ----

func newTestClient(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	cli := NewClient(srv.URL, "test-token", WithRateLimiting(false))
	return srv, cli
}

// writeJSONEnvelope writes the standard API envelope with data and optional pagination.
func writeJSONEnvelope(w http.ResponseWriter, data interface{}, pag *types.Pagination) {
	type envelope struct {
		Data       interface{}       `json:"data,omitempty"`
		Pagination *types.Pagination `json:"pagination,omitempty"`
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(envelope{Data: data, Pagination: pag})
}

// writeErrorEnvelope writes a non-2xx API error response.
func writeErrorEnvelope(w http.ResponseWriter, statusCode int, apiErrors []types.APIError) {
	type envelope struct {
		Errors []types.APIError `json:"errors"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(envelope{Errors: apiErrors})
}

// ---- NewClient / options ----

func TestNewClient_Defaults(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "mytoken")
	if cli.baseURL != "https://example.sentinelone.net" {
		t.Errorf("unexpected baseURL: %q", cli.baseURL)
	}
	if cli.apiToken != "mytoken" {
		t.Errorf("unexpected apiToken: %q", cli.apiToken)
	}
	if cli.maxRetries != defaultMaxRetries {
		t.Errorf("expected maxRetries=%d, got %d", defaultMaxRetries, cli.maxRetries)
	}
	if cli.httpClient == nil {
		t.Fatal("expected non-nil httpClient")
	}
	if cli.httpClient.Timeout != defaultTimeout {
		t.Errorf("expected timeout %v, got %v", defaultTimeout, cli.httpClient.Timeout)
	}
	if cli.rateLimits == nil {
		t.Error("expected rate limits to be enabled by default")
	}
	if cli.Accounts == nil || cli.Sites == nil || cli.RBAC == nil || cli.Users == nil {
		t.Error("expected all sub-clients to be initialised")
	}
}

func TestNewClient_TrailingSlashStripped(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net/", "tok")
	if strings.HasSuffix(cli.baseURL, "/") {
		t.Errorf("trailing slash not stripped: %q", cli.baseURL)
	}
}

func TestNewClient_MultipleTrailingSlashes(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net///", "tok")
	if strings.HasSuffix(cli.baseURL, "/") {
		t.Errorf("trailing slash not stripped: %q", cli.baseURL)
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 42 * time.Second}
	cli := NewClient("https://example.sentinelone.net", "tok", WithHTTPClient(custom))
	if cli.httpClient != custom {
		t.Error("expected custom http client to be used")
	}
}

func TestNewClient_WithTimeout(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithTimeout(99*time.Second))
	if cli.httpClient.Timeout != 99*time.Second {
		t.Errorf("expected 99s timeout, got %v", cli.httpClient.Timeout)
	}
}

func TestNewClient_WithRateLimiting_Disable(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(false))
	if cli.rateLimits != nil {
		t.Error("expected nil rateLimits when disabled")
	}
}

func TestNewClient_WithRateLimiting_Enable(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok",
		WithRateLimiting(false),
		WithRateLimiting(true),
	)
	if cli.rateLimits == nil {
		t.Error("expected non-nil rateLimits when enabled")
	}
}

func TestNewClient_WithMaxRetries(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithMaxRetries(7))
	if cli.maxRetries != 7 {
		t.Errorf("expected maxRetries=7, got %d", cli.maxRetries)
	}
}

func TestNewClient_WithMaxRetries_Zero(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithMaxRetries(0))
	if cli.maxRetries != 0 {
		t.Errorf("expected maxRetries=0, got %d", cli.maxRetries)
	}
}

// ---- AsResponseError ----

func TestAsResponseError_Success(t *testing.T) {
	orig := &types.ResponseError{StatusCode: 403}
	got, ok := AsResponseError(orig)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if got != orig {
		t.Error("expected same pointer")
	}
}

func TestAsResponseError_PlainError(t *testing.T) {
	_, ok := AsResponseError(errors.New("plain"))
	if ok {
		t.Fatal("expected ok=false")
	}
}

// ---- buildURL ----

func TestBuildURL_NoParams(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok")
	got := cli.buildURL("/accounts", nil)
	want := "https://example.sentinelone.net" + apiPathPrefix + "/accounts"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildURL_WithParams(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok")
	params := url.Values{}
	params.Set("limit", "10")
	got := cli.buildURL("/accounts", params)
	if !strings.Contains(got, "?limit=10") {
		t.Errorf("missing query string in %q", got)
	}
}

func TestBuildURL_EmptyParams(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok")
	got := cli.buildURL("/accounts", url.Values{})
	if strings.Contains(got, "?") {
		t.Errorf("unexpected query string in %q", got)
	}
}

// ---- decodeRawResponse ----

func TestDecodeRawResponse_EmptyBody(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok")
	raw, err := cli.decodeRawResponse([]byte{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if raw.Pagination != nil {
		t.Error("expected nil pagination for empty body")
	}
}

func TestDecodeRawResponse_ValidBody_NoOut(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok")
	cursor := "abc"
	body, _ := json.Marshal(map[string]interface{}{
		"pagination": map[string]interface{}{"nextCursor": cursor, "totalItems": 5},
		"data":       []string{"a", "b"},
	})
	raw, err := cli.decodeRawResponse(body, nil)
	if err != nil {
		t.Fatal(err)
	}
	if raw.Pagination == nil {
		t.Fatal("expected pagination")
	}
	if *raw.Pagination.NextCursor != cursor {
		t.Errorf("unexpected cursor: %q", *raw.Pagination.NextCursor)
	}
}

func TestDecodeRawResponse_WithOut(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok")
	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{"id": "acc1", "name": "Acme"},
	})
	var account types.Account
	_, err := cli.decodeRawResponse(body, &account)
	if err != nil {
		t.Fatal(err)
	}
	if account.ID != "acc1" {
		t.Errorf("expected id=acc1, got %q", account.ID)
	}
}

func TestDecodeRawResponse_InvalidJSON(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok")
	_, err := cli.decodeRawResponse([]byte("not-json"), nil)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "decode response") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDecodeRawResponse_InvalidDataJSON(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok")
	// data is a JSON object but out expects a string
	body, _ := json.Marshal(map[string]interface{}{"data": map[string]interface{}{"key": "val"}})
	var out string
	_, err := cli.decodeRawResponse(body, &out)
	if err == nil {
		t.Fatal("expected error for mismatched data type")
	}
	if !strings.Contains(err.Error(), "decode data") {
		t.Errorf("unexpected error: %v", err)
	}
}

// ---- parseErrors ----

func TestParseErrors_Nil(t *testing.T) {
	errs := parseErrors(nil)
	if errs != nil {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestParseErrors_Empty(t *testing.T) {
	errs := parseErrors([]json.RawMessage{})
	if errs != nil {
		t.Errorf("expected nil for empty slice, got %v", errs)
	}
}

func TestParseErrors_ValidErrors(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage(`{"code": 400, "message": "bad request", "detail": "missing field"}`),
		json.RawMessage(`{"code": 401, "message": "unauthorized"}`),
	}
	errs := parseErrors(raw)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
	if errs[0].Code != 400 || errs[0].Message != "bad request" {
		t.Errorf("unexpected first error: %+v", errs[0])
	}
	if errs[1].Code != 401 {
		t.Errorf("unexpected second error: %+v", errs[1])
	}
}

func TestParseErrors_MixedValidInvalid(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage(`{"code": 400, "message": "bad request"}`),
		json.RawMessage(`not-valid-json`),
	}
	errs := parseErrors(raw)
	if len(errs) != 1 {
		t.Errorf("expected 1 valid error, got %d", len(errs))
	}
}

// ---- do() core tests ----

// failTransport causes httpClient.Do to return an error.
type failTransport struct{ err error }

func (ft failTransport) RoundTrip(*http.Request) (*http.Response, error) { return nil, ft.err }

// errBodyReader returns a read error when the response body is read.
type errBodyReader struct{}

func (errBodyReader) Read(_ []byte) (int, error) { return 0, errors.New("body read error") }
func (errBodyReader) Close() error                { return nil }

// errBodyTransport returns a valid response but with an unreadable body.
type errBodyTransport struct{}

func (errBodyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       errBodyReader{},
	}, nil
}

// unmarshalableBody cannot be JSON-marshaled (channels).
type unmarshalableBody struct {
	C chan int `json:"c"`
}

func TestDo_Success_Get(t *testing.T) {
	srv, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		auth := r.Header.Get("Authorization")
		if auth != "ApiToken test-token" {
			t.Errorf("unexpected auth header: %q", auth)
		}
		writeJSONEnvelope(w, []string{"a", "b"}, &types.Pagination{TotalItems: 2})
	})
	_ = srv

	var out []string
	pag, err := cli.get(context.Background(), "/accounts", nil, &out)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 items, got %d", len(out))
	}
	if pag.TotalItems != 2 {
		t.Errorf("unexpected total: %d", pag.TotalItems)
	}
}

func TestDo_Success_Post(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected JSON content type, got %q", r.Header.Get("Content-Type"))
		}
		writeJSONEnvelope(w, map[string]string{"id": "new1"}, nil)
	})
	var out map[string]string
	_, err := cli.post(context.Background(), "/accounts", map[string]string{"name": "test"}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if out["id"] != "new1" {
		t.Errorf("unexpected id: %q", out["id"])
	}
}

func TestDo_Success_Put(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSONEnvelope(w, map[string]string{"id": "upd1"}, nil)
	})
	var out map[string]string
	_, err := cli.put(context.Background(), "/accounts/1", map[string]string{"name": "updated"}, &out)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDo_Success_Delete(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	_, err := cli.delete(context.Background(), "/accounts/1", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDo_NilBody_NoContentTypeHeader(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			t.Errorf("expected no Content-Type for nil body, got %q", r.Header.Get("Content-Type"))
		}
		writeJSONEnvelope(w, nil, nil)
	})
	_, err := cli.get(context.Background(), "/accounts", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDo_NonSuccessStatus(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusNotFound, []types.APIError{{Code: 404, Message: "not found"}})
	})
	_, err := cli.get(context.Background(), "/accounts/999", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	respErr, ok := AsResponseError(err)
	if !ok {
		t.Fatalf("expected *types.ResponseError, got %T: %v", err, err)
	}
	if respErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", respErr.StatusCode)
	}
	if len(respErr.Errors) != 1 || respErr.Errors[0].Code != 404 {
		t.Errorf("unexpected errors: %+v", respErr.Errors)
	}
}

func TestDo_NonSuccessStatus_NoErrors(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{}"))
	})
	_, err := cli.get(context.Background(), "/accounts", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	respErr, ok := AsResponseError(err)
	if !ok {
		t.Fatalf("expected *types.ResponseError, got %T", err)
	}
	if respErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", respErr.StatusCode)
	}
}

func TestDo_429_RetrySuccess(t *testing.T) {
	attempts := 0
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("{}"))
			return
		}
		writeJSONEnvelope(w, []string{"ok"}, nil)
	})
	cli.maxRetries = 5

	var out []string
	_, err := cli.get(context.Background(), "/accounts", nil, &out)
	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestDo_429_MaxRetriesExhausted(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("{}"))
	})
	cli.maxRetries = 0 // no retries

	_, err := cli.get(context.Background(), "/accounts", nil, nil)
	if err == nil {
		t.Fatal("expected error when max retries exhausted")
	}
	respErr, ok := AsResponseError(err)
	if !ok {
		t.Fatalf("expected *types.ResponseError, got %T: %v", err, err)
	}
	if respErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", respErr.StatusCode)
	}
}

func TestDo_429_ContextCancelledDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		// Cancel the context shortly after sending the 429, giving the select
		// time to enter the wait branch.
		go func() {
			time.Sleep(20 * time.Millisecond)
			cancel()
		}()
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("{}"))
	})
	cli.maxRetries = 5

	_, err := cli.do(ctx, http.MethodGet, "/accounts", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error when context cancelled during backoff")
	}
	if !strings.Contains(err.Error(), "context done during 429 backoff") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDo_MarshalError(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(false))
	_, err := cli.do(context.Background(), http.MethodPost, "/accounts", nil,
		unmarshalableBody{C: make(chan int)}, nil)
	if err == nil {
		t.Fatal("expected marshal error")
	}
	if !strings.Contains(err.Error(), "marshal request") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDo_HTTPRequestError(t *testing.T) {
	transportErr := errors.New("connection refused")
	cli := NewClient("https://example.sentinelone.net", "tok",
		WithHTTPClient(&http.Client{Transport: failTransport{err: transportErr}}),
		WithRateLimiting(false))

	_, err := cli.get(context.Background(), "/accounts", nil, nil)
	if err == nil {
		t.Fatal("expected error from transport failure")
	}
	if !strings.Contains(err.Error(), "HTTP request") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDo_BodyReadError(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok",
		WithHTTPClient(&http.Client{Transport: errBodyTransport{}}),
		WithRateLimiting(false))

	_, err := cli.get(context.Background(), "/accounts", nil, nil)
	if err == nil {
		t.Fatal("expected body read error")
	}
	if !strings.Contains(err.Error(), "read response") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDo_InvalidResponseJSON(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("this is not json"))
	})
	_, err := cli.get(context.Background(), "/accounts", nil, nil)
	if err == nil {
		t.Fatal("expected JSON decode error")
	}
	if !strings.Contains(err.Error(), "decode response") {
		t.Errorf("unexpected error: %v", err)
	}
}

// ---- ListParams.values() ----

func TestListParams_Values_NilReceiver(t *testing.T) {
	var p *ListParams
	vals := p.values()
	if len(vals) != 0 {
		t.Errorf("expected empty values for nil receiver, got %v", vals)
	}
}

func TestListParams_Values_Empty(t *testing.T) {
	p := &ListParams{}
	vals := p.values()
	if len(vals) != 0 {
		t.Errorf("expected empty values, got %v", vals)
	}
}

func TestListParams_Values_AllFields(t *testing.T) {
	p := &ListParams{
		Skip:      IntPtr(50),
		Limit:     IntPtr(100),
		Cursor:    StringPtr("cursor123"),
		CountOnly: BoolPtr(true),
		SkipCount: BoolPtr(false),
		SortBy:    StringPtr("name"),
		SortOrder: StringPtr("desc"),
	}
	vals := p.values()

	checks := map[string]string{
		"skip":      "50",
		"limit":     "100",
		"cursor":    "cursor123",
		"countOnly": "true",
		"skipCount": "false",
		"sortBy":    "name",
		"sortOrder": "desc",
	}
	for key, want := range checks {
		if got := vals.Get(key); got != want {
			t.Errorf("param %q: got %q, want %q", key, got, want)
		}
	}
}

// ---- setBool / setString / setStringSlice ----

func TestSetBool_Nil(t *testing.T) {
	v := url.Values{}
	setBool(v, "key", nil)
	if v.Get("key") != "" {
		t.Errorf("expected empty for nil bool")
	}
}

func TestSetBool_True(t *testing.T) {
	v := url.Values{}
	setBool(v, "key", BoolPtr(true))
	if v.Get("key") != "true" {
		t.Errorf("expected true, got %q", v.Get("key"))
	}
}

func TestSetBool_False(t *testing.T) {
	v := url.Values{}
	setBool(v, "key", BoolPtr(false))
	if v.Get("key") != "false" {
		t.Errorf("expected false, got %q", v.Get("key"))
	}
}

func TestSetString_Nil(t *testing.T) {
	v := url.Values{}
	setString(v, "key", nil)
	if v.Get("key") != "" {
		t.Errorf("expected empty for nil string")
	}
}

func TestSetString_Empty(t *testing.T) {
	v := url.Values{}
	empty := ""
	setString(v, "key", &empty)
	if v.Has("key") {
		t.Errorf("expected key absent for empty string")
	}
}

func TestSetString_NonEmpty(t *testing.T) {
	v := url.Values{}
	setString(v, "key", StringPtr("value"))
	if v.Get("key") != "value" {
		t.Errorf("expected value, got %q", v.Get("key"))
	}
}

func TestSetStringSlice_Empty(t *testing.T) {
	v := url.Values{}
	setStringSlice(v, "ids", []string{})
	if v.Has("ids") {
		t.Errorf("expected key absent for empty slice")
	}
}

func TestSetStringSlice_NonEmpty(t *testing.T) {
	v := url.Values{}
	setStringSlice(v, "ids", []string{"a", "b", "c"})
	if v.Get("ids") != "a,b,c" {
		t.Errorf("expected comma-joined, got %q", v.Get("ids"))
	}
}

// ---- BoolPtr / IntPtr / StringPtr ----

func TestBoolPtr_True(t *testing.T) {
	p := BoolPtr(true)
	if p == nil || !*p {
		t.Error("expected *true")
	}
}

func TestBoolPtr_False(t *testing.T) {
	p := BoolPtr(false)
	if p == nil || *p {
		t.Error("expected *false")
	}
}

func TestIntPtr(t *testing.T) {
	p := IntPtr(42)
	if p == nil || *p != 42 {
		t.Error("expected *42")
	}
}

func TestStringPtr(t *testing.T) {
	p := StringPtr("hello")
	if p == nil || *p != "hello" {
		t.Error("expected *hello")
	}
}

// ---- request with query params ----

func TestDo_GetWithParams(t *testing.T) {
	var receivedQuery url.Values
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.Query()
		writeJSONEnvelope(w, nil, nil)
	})
	params := url.Values{}
	params.Set("limit", "5")
	params.Set("skip", "10")
	_, err := cli.get(context.Background(), "/accounts", params, nil)
	if err != nil {
		t.Fatal(err)
	}
	if receivedQuery.Get("limit") != "5" || receivedQuery.Get("skip") != "10" {
		t.Errorf("unexpected query params: %v", receivedQuery)
	}
}

// ---- io.NopCloser for edge cases ----

func TestDo_EmptyResponseBody(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// write nothing
	})
	_, err := cli.get(context.Background(), "/accounts", nil, nil)
	if err != nil {
		t.Fatalf("expected success for empty body, got: %v", err)
	}
}

// Verify NopCloser body closes without error (indirectly through a successful request)
func TestDo_PaginationPassthrough(t *testing.T) {
	cursor := "cursor-xyz"
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSONEnvelope(w, []int{1, 2}, &types.Pagination{
			NextCursor: &cursor,
			TotalItems: 99,
		})
	})
	var out []int
	pag, err := cli.get(context.Background(), "/accounts", nil, &out)
	if err != nil {
		t.Fatal(err)
	}
	if pag == nil || pag.TotalItems != 99 {
		t.Errorf("unexpected pagination: %+v", pag)
	}
	if pag.NextCursor == nil || *pag.NextCursor != cursor {
		t.Errorf("unexpected cursor: %v", pag.NextCursor)
	}
}

func TestDo_InvalidMethod(t *testing.T) {
	// http.NewRequestWithContext rejects methods containing invalid characters.
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(false))
	_, err := cli.do(context.Background(), "IN\nVALID", "/accounts", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for invalid HTTP method")
	}
	if !strings.Contains(err.Error(), "build request") {
		t.Errorf("unexpected error: %v", err)
	}
}

// ensure io.ReadAll is covered (via successful response with actual body)
func TestDo_PostNilOut(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"data": null}`)
	})
	_, err := cli.post(context.Background(), "/users/logout", struct{}{}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
