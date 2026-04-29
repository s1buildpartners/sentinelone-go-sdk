package sentinelone

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ---- buildDefaultRateLimits ----

func TestBuildDefaultRateLimits_NonEmpty(t *testing.T) {
	limits := buildDefaultRateLimits()
	if len(limits) == 0 {
		t.Fatal("expected non-empty rate limits")
	}
}

func TestBuildDefaultRateLimits_SortedLongestFirst(t *testing.T) {
	limits := buildDefaultRateLimits()
	for i := 1; i < len(limits); i++ {
		if len(limits[i-1].prefix) < len(limits[i].prefix) {
			t.Errorf("limits not sorted longest-first at index %d: %q (len %d) before %q (len %d)",
				i, limits[i-1].prefix, len(limits[i-1].prefix),
				limits[i].prefix, len(limits[i].prefix))
		}
	}
}

func TestBuildDefaultRateLimits_AllNonNilLimiters(t *testing.T) {
	limits := buildDefaultRateLimits()
	for _, rl := range limits {
		if rl.limiter == nil {
			t.Errorf("nil limiter for prefix %q", rl.prefix)
		}
		if rl.prefix == "" {
			t.Error("empty prefix in rate limit rules")
		}
	}
}

// The /assets path has rps=0, meaning unlimited (rate.Inf).
func TestBuildDefaultRateLimits_UnlimitedForAssets(t *testing.T) {
	limits := buildDefaultRateLimits()
	for _, rl := range limits {
		if rl.prefix == "/assets" {
			for i := 0; i < 20; i++ {
				if !rl.limiter.Allow() {
					t.Error("expected unlimited limiter to always allow")
					break
				}
			}
			return
		}
	}
	t.Error("/assets prefix not found in rate limits")
}

// ---- limiterFor ----

func TestLimiterFor_ExactMatch(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(true))
	lim := cli.limiterFor("/accounts")
	if lim == nil {
		t.Fatal("expected limiter for /accounts")
	}
}

func TestLimiterFor_PrefixMatch(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(true))
	lim := cli.limiterFor("/accounts/123")
	if lim == nil {
		t.Fatal("expected limiter for /accounts/123 (prefix /accounts)")
	}
}

func TestLimiterFor_LongerPrefixWins(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(true))
	// /threats/ (longer prefix) should match before /threats
	limSlash := cli.limiterFor("/threats/123/notes")
	limNoSlash := cli.limiterFor("/threats")
	if limSlash == limNoSlash {
		t.Error("expected /threats/123/notes to use /threats/ limiter, not /threats")
	}
}

func TestLimiterFor_UsersLogin_SpecificMatch(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(true))
	loginLim := cli.limiterFor("/users/login")
	usersLim := cli.limiterFor("/users/profile")
	if loginLim == usersLim {
		t.Error("expected /users/login to use dedicated limiter, not /users")
	}
}

func TestLimiterFor_NoMatch(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(true))
	lim := cli.limiterFor("/totally-unknown-path")
	if lim != nil {
		t.Errorf("expected nil limiter for unknown path, got %v", lim)
	}
}

func TestLimiterFor_Disabled(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(false))
	lim := cli.limiterFor("/accounts")
	if lim != nil {
		t.Error("expected nil limiter when rate limiting disabled")
	}
}

func TestLimiterFor_DVPathOrder(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(true))
	pingLim := cli.limiterFor("/dv/events/pq-ping")
	pqLim := cli.limiterFor("/dv/events/pq")
	if pingLim == pqLim {
		t.Error("expected /dv/events/pq-ping to use its own limiter, not /dv/events/pq")
	}
}

func TestLimiterFor_SystemStatusSubpaths(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(true))
	cacheLim := cli.limiterFor("/system/status/cache")
	statusLim := cli.limiterFor("/system/status")
	if cacheLim == statusLim {
		t.Error("expected /system/status/cache to use its own limiter")
	}
}

func TestLimiterFor_HyperAutomateWebhookBeforeAPI(t *testing.T) {
	cli := NewClient("https://example.sentinelone.net", "tok", WithRateLimiting(true))
	webhookLim := cli.limiterFor("/hyper-automate/webhook/123")
	apiLim := cli.limiterFor("/hyper-automate/api/call")
	if webhookLim == apiLim {
		t.Error("expected /hyper-automate/webhook to use its own limiter")
	}
}

// ---- retryAfterDuration ----

func TestRetryAfterDuration_NoHeader(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	d := retryAfterDuration(resp)
	if d != 5*time.Second {
		t.Errorf("expected 5s fallback, got %v", d)
	}
}

func TestRetryAfterDuration_EmptyHeader(t *testing.T) {
	resp := &http.Response{Header: http.Header{"Retry-After": {""}}}
	d := retryAfterDuration(resp)
	if d != 5*time.Second {
		t.Errorf("expected 5s fallback for empty header, got %v", d)
	}
}

func TestRetryAfterDuration_ValidSeconds(t *testing.T) {
	resp := &http.Response{Header: http.Header{"Retry-After": {"10"}}}
	d := retryAfterDuration(resp)
	if d != 10*time.Second {
		t.Errorf("expected 10s, got %v", d)
	}
}

func TestRetryAfterDuration_ValidSeconds_Large(t *testing.T) {
	resp := &http.Response{Header: http.Header{"Retry-After": {"120"}}}
	d := retryAfterDuration(resp)
	if d != 120*time.Second {
		t.Errorf("expected 120s, got %v", d)
	}
}

func TestRetryAfterDuration_InvalidValue(t *testing.T) {
	resp := &http.Response{Header: http.Header{"Retry-After": {"not-a-number"}}}
	d := retryAfterDuration(resp)
	if d != 5*time.Second {
		t.Errorf("expected 5s fallback for invalid header, got %v", d)
	}
}

func TestRetryAfterDuration_ZeroValue(t *testing.T) {
	resp := &http.Response{Header: http.Header{"Retry-After": {"0"}}}
	d := retryAfterDuration(resp)
	if d != 5*time.Second {
		t.Errorf("expected 5s fallback for 0, got %v", d)
	}
}

func TestRetryAfterDuration_NegativeValue(t *testing.T) {
	resp := &http.Response{Header: http.Header{"Retry-After": {"-5"}}}
	d := retryAfterDuration(resp)
	if d != 5*time.Second {
		t.Errorf("expected 5s fallback for negative value, got %v", d)
	}
}

// ---- rate limit wait cancellation integration ----

func TestDo_RateLimitWaitCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	t.Cleanup(srv.Close)

	cli := NewClient(srv.URL, "tok", WithRateLimiting(true), WithMaxRetries(0))

	// Exhaust the burst token for /users/login (burst=1, rate≈1 req/s).
	_, _ = cli.do(context.Background(), http.MethodGet, "/users/login", nil, nil, nil)

	// Pre-cancel the context so lim.Wait cannot wait for the next token.
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := cli.do(cancelledCtx, http.MethodGet, "/users/login", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error when context cancelled before rate limit grants token")
	}
	if !strings.Contains(err.Error(), "rate limit wait") {
		t.Errorf("expected rate limit wait error, got: %v", err)
	}
}
