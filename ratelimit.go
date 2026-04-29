package sentinelone

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// pathRateLimit pairs a path prefix with its token-bucket limiter.
// Rules are stored sorted longest-prefix-first so that the first matching
// rule in a linear scan is always the most specific one.
type pathRateLimit struct {
	prefix  string
	limiter *rate.Limiter
}

const (
	secsPerHour         = 3600 // seconds in one hour
	scanHoursPerRequest = 8    // /application-management/scan: 1 request every N hours
)

// buildDefaultRateLimits constructs the per-path token-bucket limiters from
// the official SentinelOne MGMT API rate-limit table.  Each client gets its
// own set of limiters so that clients backed by different API tokens do not
// share state.
func buildDefaultRateLimits() []pathRateLimit { //nolint:funlen
	type rule struct {
		prefix string
		rps    float64 // requests per second; 0 = rate.Inf (unlimited)
		burst  int     // token-bucket burst capacity
	}

	rules := []rule{
		// Path prefixes are relative to /web/api/v2.1 (apiPathPrefix).
		// Burst=1 when the official table lists "none".
		{"/accounts", 100, 10},
		{"/activities", 10, 30},
		{"/agents", 25, 1},
		// scan: 1 request every scanHoursPerRequest hours
		{"/application-management/scan", 1.0 / (scanHoursPerRequest * secsPerHour), 1},
		{"/application-management/inventory", 3, 1},
		// assets: 10 000 req/s per Console — treated as unlimited
		{"/assets", 0, 1},
		// cloud-detection and detection-library: 30 req/min
		{"/cloud-detection/rules", 30.0 / 60, 1},
		{"/cloud-funnel", 2, 10},
		{"/cloudonboarding", 30, 100},
		{"/cloudsecurity", 25, 100},
		{"/cnapp", 25, 100},
		{"/detection-library/rules", 30.0 / 60, 1},
		// dv/events/pq-ping must come before dv/events/pq (longer prefix)
		{"/dv/events/pq-ping", 4, 1},
		{"/dv/events/pq", 4, 1},
		// dv/init-query: 1 req/min
		{"/dv/init-query", 1.0 / 60, 1},
		{"/dv/query-status", 1, 1},
		{"/exclusions", 25, 25},
		{"/export", 1, 30},
		// hyper-automate/webhook must come before hyper-automate/api (longer)
		{"/hyper-automate/webhook", 150, 1},
		{"/hyper-automate/api", 60, 1},
		{"/identity", 50, 100},
		{"/ranger-ad", 50, 100},
		{"/rbac", 50, 100},
		{"/remote-ops", 1, 5},
		{"/remote-scripts", 1, 30},
		{"/service-users", 5, 10},
		{"/singularity-marketplace", 5, 20},
		{"/sites", 100, 10},
		// system/status sub-paths must come before /system/status (longer)
		{"/system/status/cache", 2, 1},
		{"/system/status/db", 2, 1},
		{"/system/status", 2, 1},
		// tag-manager: docs list "1 request per API token" with no time unit;
		// treating conservatively as 1 req/s.
		{"/tag-manager", 1, 50},
		// threat-intelligence: 30 req/min (dedicated Consoles only)
		{"/threat-intelligence", 30.0 / 60, 1},
		// /threats/ (trailing slash) matches sub-paths like /threats/{id}/notes
		// and must sort before /threats so the longer prefix wins.
		{"/threats/", 100, 1000},
		{"/threats", 10, 50},
		// update/agent/download: 2 req/min
		{"/update/agent/download", 2.0 / 60, 1},
		{"/upgrade-policy", 1, 50},
		{"/upload", 1, 5},
		// /user-groups, /users/login, /users, /user — ordering matters here:
		// longer prefixes must appear first so /user does not swallow the others.
		{"/user-groups", 15, 50},
		{"/users/login", 1, 1},
		{"/users", 40, 80},
		{"/user", 2, 5},
		{"/xdr", 25, 1},
	}

	limits := make([]pathRateLimit, 0, len(rules))

	for _, entry := range rules {
		var lim *rate.Limiter

		if entry.rps == 0 {
			lim = rate.NewLimiter(rate.Inf, entry.burst)
		} else {
			lim = rate.NewLimiter(rate.Limit(entry.rps), entry.burst)
		}

		limits = append(limits, pathRateLimit{prefix: entry.prefix, limiter: lim})
	}

	// Sort longest-prefix-first so the linear scan in limiterFor always picks
	// the most specific matching rule.
	sort.Slice(limits, func(i, j int) bool {
		return len(limits[i].prefix) > len(limits[j].prefix)
	})

	return limits
}

// limiterFor returns the rate.Limiter for the given API path, or nil when rate
// limiting is disabled (c.rateLimits is empty).  Matching is longest-prefix-first.
func (c *Client) limiterFor(path string) *rate.Limiter {
	for _, rl := range c.rateLimits {
		if strings.HasPrefix(path, rl.prefix) {
			return rl.limiter
		}
	}

	return nil
}

// retryAfterDuration reads the Retry-After response header and converts it to
// a wait duration.  Falls back to 5 seconds when the header is absent or
// cannot be parsed.
func retryAfterDuration(resp *http.Response) time.Duration {
	const fallback = 5 * time.Second

	h := resp.Header.Get("Retry-After")
	if h == "" {
		return fallback
	}

	secs, err := strconv.Atoi(h)
	if err == nil && secs > 0 {
		return time.Duration(secs) * time.Second
	}

	return fallback
}
