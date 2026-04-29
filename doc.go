// Package sentinelone provides a Go client for the SentinelOne Management API v2.1.
//
// # Overview
//
// The client covers four API groups: Accounts, Sites, RBAC (roles), and Users.
// Every call requires a context, returns typed structs, and maps non-2xx responses
// to a [ResponseError] so callers can inspect the HTTP status code and any API
// error messages in one place.
//
// # Authentication
//
// All requests authenticate via an API token sent in the Authorization header
// as "ApiToken <token>".  Generate a token in the SentinelOne console under
// My User → Actions → API Token Operations, or programmatically via
// [Client.GenerateAPIToken].
//
// # Creating a client
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
// # Error handling
//
// Network errors and JSON decode failures are returned as ordinary Go errors.
// HTTP 4xx/5xx responses are returned as *[ResponseError], which carries the
// status code and the list of [APIError] values from the response body:
//
//	accounts, _, err := client.ListAccounts(ctx, nil)
//	if err != nil {
//	    var respErr *sentinelone.ResponseError
//	    if errors.As(err, &respErr) {
//	        fmt.Printf("API error %d: %v\n", respErr.StatusCode, respErr.Errors)
//	    }
//	    return err
//	}
//
// # Pagination
//
// List endpoints return a [Pagination] value alongside the result slice.
// Use cursor-based iteration for result sets larger than 1,000 items:
//
//	var cursor *string
//	for {
//	    accounts, pag, err := client.ListAccounts(ctx, &sentinelone.ListAccountsParams{
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
//	accounts, _, err := client.ListAccounts(ctx, &sentinelone.ListAccountsParams{
//	    State: "active",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if len(accounts) == 0 {
//	    return
//	}
//
//	updated, err := client.UpdateAccount(ctx, accounts[0].ID, sentinelone.UpdateAccountRequest{
//	    Data: sentinelone.UpdateAccountData{Name: "Renamed Account"},
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Updated:", updated.Name)
//
// Reactivate an expired account with a new expiry date:
//
//	_, err = client.ReactivateAccount(ctx, accountID, sentinelone.ReactivateAccountRequest{
//	    Data: sentinelone.ReactivateAccountData{
//	        Expiration: sentinelone.StringPtr("2027-01-01T00:00:00Z"),
//	    },
//	})
//
// # Sites
//
// List all sites in an account, then retrieve a single site by ID:
//
//	resp, _, err := client.ListSites(ctx, &sentinelone.ListSitesParams{
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
//	site, err := client.CreateSite(ctx, sentinelone.CreateSiteRequest{
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
//	token, err := client.GetSiteToken(ctx, site.ID)
//	fmt.Println("Registration token:", token.Token)
//
// # RBAC
//
// List all custom (non-predefined) roles visible in an account, then fetch
// the full permission details for one of them:
//
//	roles, _, err := client.ListRoles(ctx, &sentinelone.ListRolesParams{
//	    AccountIDs:     []string{"225494730938493804"},
//	    PredefinedRole: sentinelone.BoolPtr(false),
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	role, err := client.GetRole(ctx, roles[0].ID, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, page := range role.Pages {
//	    fmt.Println(page.Name)
//	}
//
// Create a role scoped to a specific account:
//
//	newRole, err := client.CreateRole(ctx, sentinelone.CreateRoleRequest{
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
//	users, _, err := client.ListUsers(ctx, &sentinelone.ListUsersParams{
//	    TwoFAEnabled: sentinelone.BoolPtr(false),
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	ids := make([]string, len(users))
//	for i, u := range users {
//	    ids[i] = u.ID
//	}
//	if _, err := client.Enroll2FA(ctx, sentinelone.UserIDsRequest{
//	    Data: sentinelone.UserIDsData{UserIDs: ids},
//	}); err != nil {
//	    log.Fatal(err)
//	}
//
// Create a new user and assign them a role in a specific site:
//
//	user, err := client.CreateUser(ctx, sentinelone.CreateUserRequest{
//	    Data: sentinelone.CreateUserData{
//	        Email:    "alice@example.com",
//	        FullName: "Alice Example",
//	        Scope:    "site",
//	        ScopeRoles: []sentinelone.UserScopeRole{
//	            {ID: "225494730938493805", Name: "Production East", AccountName: "Acme", RoleID: "225494730938493900"},
//	        },
//	    },
//	})
//
// Log in with credentials and store the returned token for subsequent calls:
//
//	resp, err := client.Login(ctx, sentinelone.LoginRequest{
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
