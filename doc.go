// Package sentinelone provides a Go client for the SentinelOne Management API v2.1.
//
// # Overview
//
// API calls are organised into four sub-clients, each accessed as a field on
// the root [Client]:
//
//   - [Client.Accounts] — account lifecycle, policy, and uninstall passwords
//   - [Client.Sites]    — site lifecycle, policy, tokens, and bulk operations
//   - [Client.RBAC]     — role listing, templates, and CRUD
//   - [Client.Users]    — user CRUD, auth, 2FA, password, and API token management
//
// Every method returns typed structs from the
// [github.com/s1buildpartners/sentinelone-go-sdk/types] subpackage and maps
// non-2xx responses to a [types.ResponseError] so callers can inspect the HTTP
// status code and any API error messages in one place.
//
// # Packages
//
// Callers typically import both this package and the types subpackage:
//
//	import (
//	    s1      "github.com/s1buildpartners/sentinelone-go-sdk"
//	    s1types "github.com/s1buildpartners/sentinelone-go-sdk/types"
//	)
//
// Input types (request bodies, filter params) live in the root package.
// Domain/model types (API response structs) live in the types subpackage.
//
// # Authentication
//
// All requests authenticate via an API token sent in the Authorization header
// as "ApiToken <token>".  Generate a token in the SentinelOne console under
// My User → Actions → API Token Operations, or programmatically via
// [UsersClient.GenerateAPIToken].
//
// # Creating a client
//
// Pass credentials directly:
//
//	client := sentinelone.NewClient(
//	    "https://your-tenant.sentinelone.net",
//	    "your-api-token",
//	)
//
// Override the default 30-second timeout:
//
//	client := sentinelone.NewClient(
//	    "https://your-tenant.sentinelone.net",
//	    "your-api-token",
//	    sentinelone.WithTimeout(60*time.Second),
//	)
//
// Bring your own [net/http.Client] (useful for proxies, mTLS, or custom
// transport settings):
//
//	transport := &http.Transport{TLSClientConfig: tlsCfg}
//	client := sentinelone.NewClient(
//	    "https://your-tenant.sentinelone.net",
//	    "your-api-token",
//	    sentinelone.WithHTTPClient(&http.Client{Transport: transport}),
//	)
//
// # Credential configuration
//
// Instead of supplying credentials directly to [NewClient] you can load them
// from environment variables or a credentials file.
//
// # Loading from environment variables
//
// Set SENTINELONE_URL to the management console base URL and
// SENTINELONE_TOKEN to a valid API token, then call [NewClientFromEnv].
// An error is returned if either variable is absent or empty.
//
//	client, err := sentinelone.NewClientFromEnv()
//
// # Loading from a credentials file
//
// [NewClientFromConfig] reads credentials from an INI-style file under the
// named profile.  Use [WithProfile] to choose a profile; the "default" profile
// is used when it is omitted.  Use [WithConfigFile] to bypass the
// SENTINELONE_CONFIG environment variable and the platform default path.
//
// The default file path is platform-specific:
//   - Linux/BSD:  $XDG_CONFIG_HOME/sentinelone/credentials  (or ~/.config/…)
//   - macOS:      ~/Library/Application Support/sentinelone/credentials
//   - Windows:    %AppData%\SentinelOne\credentials
//
// File format:
//
//	# lines starting with '#' or ';' are comments
//	[default]
//	url   = https://tenant.sentinelone.net
//	token = your-api-token
//
//	[production]
//	url   = https://prod.sentinelone.net
//	token = prod-api-token
//
// Both '=' and ':' are accepted as key-value separators.
//
// Load the default profile:
//
//	client, err := sentinelone.NewClientFromConfig()
//
// Load a named profile:
//
//	client, err := sentinelone.NewClientFromConfig(sentinelone.WithProfile("production"))
//
// Specify a custom credentials file:
//
//	client, err := sentinelone.NewClientFromConfig(sentinelone.WithConfigFile("/etc/myapp/creds"))
//
// # Layered credential lookup
//
// [NewClientFromProfile] is the recommended constructor for applications that
// need to run in both CI/container environments (env vars) and on developer
// workstations (credentials file) without code changes.
//
// Priority order:
//  1. SENTINELONE_URL and SENTINELONE_TOKEN — used directly when both are set.
//  2. Credentials file — the named profile is loaded.  When [WithProfile] is
//     not provided, SENTINELONE_PROFILE is checked and then "default" is used
//     as a fallback.
//
// Examples:
//
//	// Env vars win in CI; falls back to the "default" profile locally.
//	client, err := sentinelone.NewClientFromProfile()
//
//	// Env vars win in CI; falls back to the "production" profile locally.
//	client, err := sentinelone.NewClientFromProfile(sentinelone.WithProfile("production"))
//
// # Rate Limiting
//
// The client enforces SentinelOne's published per-API-token rate limits using a
// per-path token-bucket limiter (golang.org/x/time/rate).  Rate limiting is
// enabled by default — no configuration is required.
//
// Before each request the client acquires a token for the matching path prefix.
// When the bucket is empty the call blocks until a token is available, keeping
// the sustained throughput under the API's limit.  Each [Client] instance
// maintains independent limiter state, so multiple clients backed by different
// API tokens do not share quota.
//
// If the API still returns a 429 Too Many Requests response (e.g. due to burst
// traffic from another process sharing the same token), the client reads the
// Retry-After header, waits the indicated number of seconds (defaulting to 5 s
// when the header is absent), then retries automatically — up to 3 times by
// default.  Both the proactive wait and the 429 backoff respect the
// [context.Context] passed to the method: a cancelled or timed-out context
// aborts the wait and returns the context's error immediately.
//
// To disable the built-in limiter when managing throttling externally:
//
//	client := sentinelone.NewClient(baseURL, token,
//	    sentinelone.WithRateLimiting(false),
//	)
//
// To change the number of automatic 429 retries:
//
//	client := sentinelone.NewClient(baseURL, token,
//	    sentinelone.WithMaxRetries(5), // retry up to 5 times
//	)
//
//	client := sentinelone.NewClient(baseURL, token,
//	    sentinelone.WithMaxRetries(0), // treat 429 as an error, no retry
//	)
//
// # Context
//
// Every method accepts a [context.Context] as its first argument.  The context
// controls two things:
//
//   - Cancellation: cancelling the context (e.g. via context.WithCancel or a
//     request-scoped context from an HTTP server) aborts the in-flight HTTP
//     request and causes the method to return immediately with the context's
//     error.
//
//   - Deadlines and timeouts: a deadline set on the context (context.WithDeadline
//     or context.WithTimeout) caps the total time allowed for the call,
//     regardless of the client-level WithTimeout setting.  The stricter of the
//     two limits wins.
//
// For one-off scripts or tests, context.Background() is sufficient.  In
// server or pipeline code, thread the incoming request context through so
// that cancellation and tracing propagate correctly:
//
//	// Abort the call if it takes longer than 5 seconds.
//	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
//	defer cancel()
//
//	accounts, _, err := client.Accounts.List(ctx, nil)
//
// # Error handling
//
// Network errors and JSON decode failures are returned as ordinary Go errors.
// HTTP 4xx/5xx responses are returned as *[types.ResponseError], which carries
// the status code and the list of [types.APIError] values from the response
// body.  Use [AsResponseError] to unwrap without importing the errors package:
//
//	accounts, _, err := client.Accounts.List(ctx, nil)
//	if err != nil {
//	    if respErr, ok := sentinelone.AsResponseError(err); ok {
//	        fmt.Printf("API error %d: %v\n", respErr.StatusCode, respErr.Errors)
//	    }
//	    return err
//	}
//
// # Pagination
//
// List endpoints return a [types.Pagination] value alongside the result slice.
// Use cursor-based iteration for result sets larger than 1,000 items:
//
//	var cursor *string
//	for {
//	    accounts, pag, err := client.Accounts.List(ctx, &sentinelone.ListAccountsParams{
//	        ListParams: sentinelone.ListParams{
//	            Limit:  sentinelone.IntPtr(1000),
//	            Cursor: cursor,
//	        },
//	    })
//	    if err != nil {
//	        return err
//	    }
//	    process(accounts)
//	    if pag == nil || pag.NextCursor == nil {
//	        break
//	    }
//	    cursor = pag.NextCursor
//	}
//
// # Accounts
//
// List all active accounts, then update the name of the first one found:
//
//	accounts, _, err := client.Accounts.List(ctx, &sentinelone.ListAccountsParams{
//	    State: "active",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if len(accounts) == 0 {
//	    return
//	}
//
//	updated, err := client.Accounts.Update(ctx, accounts[0].ID, sentinelone.UpdateAccountRequest{
//	    Data: sentinelone.UpdateAccountData{Name: "Renamed Account"},
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Updated:", updated.Name)
//
// Reactivate an expired account with a new expiry date:
//
//	_, err = client.Accounts.Reactivate(ctx, accountID, sentinelone.ReactivateAccountRequest{
//	    Data: sentinelone.ReactivateAccountData{
//	        Expiration: sentinelone.StringPtr("2027-01-01T00:00:00Z"),
//	    },
//	})
//
// # Sites
//
// List all sites in an account, then retrieve a single site by ID:
//
//	resp, _, err := client.Sites.List(ctx, &sentinelone.ListSitesParams{
//	    AccountID: "225494730938493804",
//	    State:     "active",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, s := range resp.Sites {
//	    fmt.Printf("%s (%s)\n", s.Name, s.ID)
//	}
//
// Create a new site and retrieve its registration token:
//
//	site, err := client.Sites.Create(ctx, sentinelone.CreateSiteRequest{
//	    Data: sentinelone.CreateSiteData{
//	        Name:      "Production East",
//	        AccountID: "225494730938493804",
//	        SiteType:  "Paid",
//	        SKU:       "complete",
//	    },
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	token, err := client.Sites.GetToken(ctx, site.ID)
//	fmt.Println("Registration token:", token.Token)
//
// # RBAC
//
// List all custom (non-predefined) roles visible in an account, then fetch
// the full permission details for one of them:
//
//	roles, _, err := client.RBAC.List(ctx, &sentinelone.ListRolesParams{
//	    AccountIDs:     []string{"225494730938493804"},
//	    PredefinedRole: sentinelone.BoolPtr(false),
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	role, err := client.RBAC.Get(ctx, roles[0].ID, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, page := range role.Pages {
//	    fmt.Println(page.Name)
//	}
//
// Create a role scoped to a specific account:
//
//	newRole, err := client.RBAC.Create(ctx, sentinelone.CreateRoleRequest{
//	    Data: sentinelone.CreateRoleData{
//	        Name:        "Acme-ReadOnly",
//	        Description: "View-only role for Acme account",
//	    },
//	    Filter: sentinelone.RoleScopeFilter{
//	        AccountIDs: []string{"225494730938493804"},
//	    },
//	})
//
// # Users
//
// List users who have not yet enabled 2FA, then force-enable it for all of them:
//
//	users, _, err := client.Users.List(ctx, &sentinelone.ListUsersParams{
//	    TwoFAEnabled: sentinelone.BoolPtr(false),
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	ids := make([]string, len(users))
//	for i, u := range users {
//	    ids[i] = u.ID
//	}
//	if _, err := client.Users.Enroll2FA(ctx, sentinelone.UserIDsRequest{
//	    Data: sentinelone.UserIDsData{UserIDs: ids},
//	}); err != nil {
//	    log.Fatal(err)
//	}
//
// Create a new user and assign them a role in a specific site:
//
//	user, err := client.Users.Create(ctx, sentinelone.CreateUserRequest{
//	    Data: sentinelone.CreateUserData{
//	        Email:    "alice@example.com",
//	        FullName: "Alice Example",
//	        Scope:    "site",
//	        ScopeRoles: []s1types.UserScopeRole{
//	            {ID: "225494730938493805", Name: "Production East", AccountName: "Acme", RoleID: "225494730938493900"},
//	        },
//	    },
//	})
//
// Log in with credentials and store the returned token for subsequent calls:
//
//	resp, err := client.Users.Login(ctx, sentinelone.LoginRequest{
//	    Username: "admin@example.com",
//	    Password: "s3cur3P@ss",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use resp.Token to create an authenticated client.
//	authedClient := sentinelone.NewClient(baseURL, resp.Token)
//
// # Helper constructors
//
// Several pointer-typed fields accept *bool, *int, or *string.  The package
// exports small helpers to create these inline without a temporary variable:
//
//	sentinelone.BoolPtr(true)    // *bool
//	sentinelone.IntPtr(100)      // *int
//	sentinelone.StringPtr("asc") // *string
package sentinelone
