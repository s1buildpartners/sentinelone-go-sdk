# sentinelone-go-sdk

A Go client library for interacting with all of the SentinelOne APIs, which includes the REST APIs, GraphQL APIs and SDL/AI SIEM APIs.

This SDK currently covers the following API groups:

| Group        | File                        | Methods                                                                    |
|--------------|-----------------------------|----------------------------------------------------------------------------|
| **Accounts** | [accounts.go](accounts.go)  | List, Get, Create, Update, policy management, uninstall passwords          |
| **Sites**    | [sites.go](sites.go)        | List, Get, Create, Update, Delete, policy management, token rotation       |
| **RBAC**     | [rbac.go](rbac.go)          | List roles, Get role template, Get, Create, Update, Delete roles           |
| **Users**    | [users.go](users.go)        | List, Get, Create, Update, Delete, auth, 2FA, API tokens, SSO, onboarding  |

---

## Requirements

- Go 1.25.9 or later
- A SentinelOne management console URL and API token

---

## Installation

```bash
go get github.com/s1buildpartners/sentinelone-go-sdk
```

---

## Authentication

All API calls authenticate via an API token sent in the `Authorization` header as `ApiToken <token>`.

There are two ways to obtain a token:

**From the console:** Log in → click your avatar (top-right) → *My User* → *Actions* → *API Token Operations* → *Generate*.

**Programmatically** (if you already have a session token):

```go
resp, err := client.Users.GenerateAPIToken(ctx, sentinelone.GenerateAPITokenRequest{})
if err != nil {
    log.Fatal(err)
}
fmt.Println("API token:", resp.Token)
```

---

## Creating a client

```go
import "github.com/s1buildpartners/sentinelone-go-sdk"

client := sentinelone.NewClient(
    "https://your-tenant.sentinelone.net",
    "your-api-token",
)
```

### Options

| Option                              | Default | Description                                                       |
|-------------------------------------|---------|-------------------------------------------------------------------|
| `sentinelone.WithTimeout(d)`        | 30 s    | Override the per-request HTTP timeout                             |
| `sentinelone.WithHTTPClient(hc)`    | —       | Provide a custom `*http.Client` (proxy, mTLS, custom transport)   |
| `sentinelone.WithRateLimiting(b)`   | `true`  | Enable or disable the built-in per-path token-bucket rate limiter |
| `sentinelone.WithMaxRetries(n)`     | `3`     | Maximum number of automatic retries on a 429 response             |

```go
import (
    "net/http"
    "time"
    "github.com/s1buildpartners/sentinelone-go-sdk"
)

// Custom timeout
client := sentinelone.NewClient(baseURL, token,
    sentinelone.WithTimeout(60*time.Second),
)

// Custom HTTP client (e.g. with a proxy)
transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
client := sentinelone.NewClient(baseURL, token,
    sentinelone.WithHTTPClient(&http.Client{Transport: transport}),
)
```

---

## Credential configuration

Rather than embedding credentials in code you can load them from environment variables or a credentials file. All three constructors accept the same `ClientOption` values as `NewClient`.

### Environment variables

| Variable | Description |
| --- | --- |
| `SENTINELONE_URL` | Management console base URL |
| `SENTINELONE_TOKEN` | API token |
| `SENTINELONE_CONFIG` | Override the credentials file path |
| `SENTINELONE_PROFILE` | Default profile name (used when profile arg is `""`) |

```go
// Both SENTINELONE_URL and SENTINELONE_TOKEN must be set.
client, err := sentinelone.NewClientFromEnv()
if err != nil {
    log.Fatal(err)
}
```

### Credentials file

`NewClientFromConfig` reads an INI-style file. The default location is platform-specific:

| Platform | Default path |
| --- | --- |
| Linux / BSD | `$XDG_CONFIG_HOME/sentinelone/credentials` (or `~/.config/sentinelone/credentials`) |
| macOS | `~/Library/Application Support/sentinelone/credentials` |
| Windows | `%AppData%\SentinelOne\credentials` |

Override the path with `SENTINELONE_CONFIG`.

**File format:**

```ini
# lines beginning with '#' or ';' are comments
[default]
url   = https://tenant.sentinelone.net
token = your-api-token

[production]
url   = https://prod.sentinelone.net
token = prod-api-token

[staging]
url   = https://staging.sentinelone.net
token = staging-api-token
```

Both `=` and `:` are accepted as key-value separators.

```go
// Load the "default" profile (pass "" for the default).
client, err := sentinelone.NewClientFromConfig("")
if err != nil {
    log.Fatal(err)
}

// Load a named profile.
client, err = sentinelone.NewClientFromConfig("production")
if err != nil {
    log.Fatal(err)
}
```

### Layered lookup (recommended)

`NewClientFromProfile` is the recommended constructor for applications that need to run in both CI/container environments (env vars) and on developer workstations (credentials file) without code changes.

**Priority order:**

1. `SENTINELONE_URL` + `SENTINELONE_TOKEN` — used directly when both are set.
2. Credentials file — the named profile is loaded. When the profile argument is `""`, `SENTINELONE_PROFILE` is checked and then `"default"` is used as a final fallback.

```go
// Env vars win in CI; falls back to the "default" profile locally.
client, err := sentinelone.NewClientFromProfile("")
if err != nil {
    log.Fatal(err)
}

// Env vars win in CI; falls back to the "production" profile locally.
client, err = sentinelone.NewClientFromProfile("production")
if err != nil {
    log.Fatal(err)
}

// Options work the same way as with NewClient.
client, err = sentinelone.NewClientFromProfile("",
    sentinelone.WithTimeout(60*time.Second),
)
```

---

## Rate limiting

The client enforces SentinelOne's published per-API-token rate limits automatically using a per-path token-bucket limiter (from `golang.org/x/time/rate`). Rate limiting is **on by default** — no configuration required.

### How it works

**Proactive (token bucket):** Before each request the client acquires a token for the matching path prefix. If the bucket is empty the call blocks until a token becomes available, keeping the sustained request rate under the API's limit. Each client instance maintains independent limiter state, so multiple clients backed by different API tokens do not interfere with each other.

**Reactive (429 retry):** If the API still returns a 429 Too Many Requests response (e.g. due to burst traffic from another process sharing the same token), the client reads the `Retry-After` header, waits the indicated number of seconds (defaulting to 5 s when the header is absent), and retries automatically — up to 3 times by default.

Both mechanisms respect the `context.Context` passed to each method. A cancelled or timed-out context will abort any in-progress wait and return the context's error to the caller.

### Rate limits enforced

A subset of the built-in limits:

| Path prefix | Sustained rate | Burst |
| ----------- | -------------- | ----- |
| `/accounts` | 100 req/s | 10 |
| `/sites` | 100 req/s | 10 |
| `/users` | 40 req/s | 80 |
| `/rbac` | 50 req/s | 100 |
| `/agents` | 25 req/s | — |
| `/threats` | 10 req/s | 50 |
| `/threats/<id>/…` | 100 req/s | 1000 |
| `/user` | 2 req/s | 5 |
| `/users/login` | 1 req/s | — |

The full table covers all documented MGMT API paths. See [ratelimit.go](ratelimit.go) for the complete list.

### Disabling rate limiting

If you manage throttling externally (your own middleware, a reverse proxy, etc.) you can turn off the built-in limiter:

```go
client := sentinelone.NewClient(baseURL, token,
    sentinelone.WithRateLimiting(false),
)
```

### Configuring retry behaviour

```go
// Retry up to 5 times on 429 before giving up
client := sentinelone.NewClient(baseURL, token,
    sentinelone.WithMaxRetries(5),
)

// Disable retries entirely (treat 429 as an error)
client := sentinelone.NewClient(baseURL, token,
    sentinelone.WithMaxRetries(0),
)
```

---

## Context

Every method takes a `context.Context` as its first argument. The context serves two purposes:

**Cancellation** — cancelling the context (for example, via `context.WithCancel` or a request-scoped context from an HTTP handler) aborts the in-flight HTTP request immediately and returns the context's error to the caller.

**Deadlines and timeouts** — a deadline on the context caps how long a call may take, independent of the client-level `WithTimeout` setting. The stricter of the two limits wins.

For scripts or tests, `context.Background()` is fine. In server or pipeline code, pass the incoming request context through so cancellation and tracing propagate correctly:

```go
// Abort if the call takes longer than 5 seconds.
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

accounts, _, err := client.Accounts.List(ctx, nil)
```

---

## Error handling

Network failures and JSON decode errors are returned as standard Go errors.

HTTP 4xx/5xx responses are returned as `*sentinelone.ResponseError`, which exposes the HTTP status code and any structured API error messages. Use `sentinelone.AsResponseError` to unwrap without importing the `errors` package yourself:

```go
accounts, _, err := client.Accounts.List(ctx, nil)
if err != nil {
    if respErr, ok := sentinelone.AsResponseError(err); ok {
        fmt.Printf("API returned HTTP %d\n", respErr.StatusCode)
        for _, e := range respErr.Errors {
            fmt.Printf("  error: %s\n", e.Message)
        }
    }
    return err
}
```

---

## Pagination

List endpoints return a `*sentinelone.Pagination` value alongside the result slice. `Pagination.NextCursor` is non-nil when more pages exist; pass it as `ListParams.Cursor` on the next call.

```go
var cursor *string

for {
    accounts, pag, err := client.Accounts.List(ctx, &sentinelone.ListAccountsParams{
        ListParams: sentinelone.ListParams{
            Limit:  sentinelone.IntPtr(1000),
            Cursor: cursor,
        },
    })
    if err != nil {
        return err
    }

    for _, a := range accounts {
        fmt.Println(a.ID, a.Name)
    }

    if pag == nil || pag.NextCursor == nil {
        break
    }
    cursor = pag.NextCursor
}
```

> **Tip:** Setting `SkipCount: sentinelone.BoolPtr(true)` skips the `COUNT(*)` query on the server and speeds up large iterations.

---

## Helper functions

Several request struct fields are typed as `*bool`, `*int`, or `*string`. The SDK provides small helpers so you can write these inline without a temporary variable:

```go
sentinelone.BoolPtr(true)      // *bool
sentinelone.IntPtr(100)        // *int
sentinelone.StringPtr("asc")   // *string
```

---

## Accounts

### List accounts

```go
ctx := context.Background()

// All active accounts
accounts, pag, err := client.Accounts.List(ctx, &sentinelone.ListAccountsParams{
    State: "active",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d accounts (total: %d)\n", len(accounts), pag.TotalItems)
```

### Get a single account

```go
account, err := client.Accounts.Get(ctx, "225494730938493804")
if err != nil {
    log.Fatal(err)
}
fmt.Println(account.Name, account.State)
```

### Create an account

```go
account, err := client.Accounts.Create(ctx, sentinelone.CreateAccountRequest{
    Data: sentinelone.CreateAccountData{
        Name:        "Acme Corp",
        AccountType: "Paid",
        Expiration:  sentinelone.StringPtr("2027-01-01T00:00:00Z"),
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Created account:", account.ID)
```

### Update an account

```go
updated, err := client.Accounts.Update(ctx, account.ID, sentinelone.UpdateAccountRequest{
    Data: sentinelone.UpdateAccountData{
        Name:        "Acme Corporation",
        ExternalID:  sentinelone.StringPtr("crm-98765"),
    },
})
```

### Manage the account policy

```go
// Read the current policy
policy, err := client.Accounts.GetPolicy(ctx, account.ID)

// Override specific engine settings
_, err = client.Accounts.UpdatePolicy(ctx, account.ID, sentinelone.UpdatePolicyRequest{
    Data: sentinelone.Policy{
        MitigationMode: "protect",
        Engines: &sentinelone.PolicyEngines{
            Reputation:  "on",
            Executables: "on",
        },
    },
})

// Revert back to the global tenant policy
err = client.Accounts.RevertPolicy(ctx, account.ID)
```

### Uninstall password

```go
// Check whether a password exists
meta, err := client.Accounts.GetUninstallPasswordMetadata(ctx, account.ID)
fmt.Println("Has password:", meta.HasPassword)

// Generate a new password (returned only once — store it securely)
pw, err := client.Accounts.GenerateUninstallPassword(ctx, account.ID)
fmt.Println("Password:", pw.Password)

// Revoke the password
err = client.Accounts.RevokeUninstallPassword(ctx, account.ID)
```

### Lifecycle operations

```go
// Immediately expire an account
err = client.Accounts.ExpireNow(ctx, account.ID)

// Reactivate it with a new expiry
_, err = client.Accounts.Reactivate(ctx, account.ID, sentinelone.ReactivateAccountRequest{
    Data: sentinelone.ReactivateAccountData{
        Expiration: sentinelone.StringPtr("2028-01-01T00:00:00Z"),
    },
})
```

---

## Sites

### List sites

```go
resp, pag, err := client.Sites.List(ctx, &sentinelone.ListSitesParams{
    AccountID: "225494730938493804",
    State:     "active",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total licenses across all sites: %d\n", resp.AllSites.TotalLicenses)
for _, s := range resp.Sites {
    fmt.Printf("  %s (%s) — %d active licenses\n", s.Name, s.ID, s.ActiveLicenses)
}
```

### Get a single site

```go
site, err := client.Sites.Get(ctx, "225494730938493805")
if err != nil {
    log.Fatal(err)
}
fmt.Println(site.Name, site.SKU)
```

### Create a site

```go
site, err := client.Sites.Create(ctx, sentinelone.CreateSiteRequest{
    Data: sentinelone.CreateSiteData{
        Name:                "Production East",
        AccountID:           "225494730938493804",
        SiteType:            "Paid",
        SKU:                 "complete",
        UnlimitedExpiration: sentinelone.BoolPtr(true),
        UnlimitedLicenses:   sentinelone.BoolPtr(true),
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Created site:", site.ID)
```

### Update a site

```go
_, err = client.Sites.Update(ctx, site.ID, sentinelone.UpdateSiteRequest{
    Data: sentinelone.UpdateSiteData{
        Description: "Primary east-coast production site",
        ExternalID:  sentinelone.StringPtr("site-east-001"),
    },
})
```

### Delete a site

```go
err = client.Sites.Delete(ctx, site.ID)
```

### Registration token

```go
// Get the current token
token, err := client.Sites.GetToken(ctx, site.ID)
fmt.Println("Registration token:", token.Token)

// Rotate the key (invalidates the old token)
newKey, err := client.Sites.RegenerateKey(ctx, site.ID)
fmt.Println("New token:", newKey.Token)
```

### Duplicate a site

```go
duplicate, err := client.Sites.Duplicate(ctx, sentinelone.DuplicateSiteRequest{
    Data: sentinelone.DuplicateSiteData{
        SiteID:     site.ID,
        Name:       "Production West",
        CopyPolicy: sentinelone.BoolPtr(true),
    },
})
```

### Bulk update sites

```go
err = client.Sites.BulkUpdate(ctx, sentinelone.BulkUpdateSitesRequest{
    Filter: sentinelone.BulkUpdateSitesFilter{
        AccountIDs: []string{"225494730938493804"},
    },
    Data: sentinelone.UpdateSiteData{
        Description: "Managed by automation",
    },
})
```

### Local upgrade authorization

```go
// Allow agents to upgrade themselves until end of month
expiry := "2026-04-30T23:59:59Z"
_, err = client.Sites.UpdateLocalAuthorization(ctx, site.ID,
    sentinelone.UpdateLocalAuthorizationRequest{
        SiteAuthorization: &expiry,
    },
)

// Revoke authorization
_, err = client.Sites.UpdateLocalAuthorization(ctx, site.ID,
    sentinelone.UpdateLocalAuthorizationRequest{},
)
```

---

## RBAC

### List roles

```go
// All custom (non-system) roles in an account
roles, pag, err := client.RBAC.List(ctx, &sentinelone.ListRolesParams{
    AccountIDs:     []string{"225494730938493804"},
    PredefinedRole: sentinelone.BoolPtr(false),
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%d custom roles\n", pag.TotalItems)
for _, r := range roles {
    fmt.Printf("  %s — %d users\n", r.Name, r.UsersInRoles)
}
```

### Get a role template (for building new roles)

```go
// Fetch the permission page structure for the account scope
template, err := client.RBAC.GetTemplate(ctx, &sentinelone.GetRoleTemplateParams{
    AccountIDs: []string{"225494730938493804"},
})
if err != nil {
    log.Fatal(err)
}

// Collect identifiers for permissions you want to grant
var permIDs []string
for _, page := range template.Pages {
    for _, perm := range page.Permissions {
        if perm.Value { // default-on permissions
            permIDs = append(permIDs, perm.Identifier)
        }
    }
}
```

### Get full permissions for a specific role

```go
role, err := client.RBAC.Get(ctx, "225494730938493900", nil)
if err != nil {
    log.Fatal(err)
}
for _, page := range role.Pages {
    fmt.Println(page.Name)
    for _, perm := range page.Permissions {
        if perm.Value {
            fmt.Printf("  [x] %s\n", perm.Title)
        }
    }
}
```

### Create a role

```go
newRole, err := client.RBAC.Create(ctx, sentinelone.CreateRoleRequest{
    Data: sentinelone.CreateRoleData{
        Name:          "Acme-ReadOnly",
        Description:   "View-only access for Acme account",
        PermissionIDs: permIDs, // from GetRoleTemplate above
    },
    Filter: sentinelone.RoleScopeFilter{
        AccountIDs: []string{"225494730938493804"},
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Created role:", newRole.ID)
```

### Update a role

```go
_, err = client.RBAC.Update(ctx, newRole.ID, sentinelone.UpdateRoleRequest{
    Data: sentinelone.UpdateRoleData{
        Name:          "Acme-ReadOnly",
        Description:   "View-only access — updated",
        PermissionIDs: permIDs,
    },
})
```

### Delete a role

```go
err = client.RBAC.Delete(ctx, newRole.ID)
```

---

## Users

### List users

```go
users, pag, err := client.Users.List(ctx, &sentinelone.ListUsersParams{
    AccountIDs:   []string{"225494730938493804"},
    EmailVerified: sentinelone.BoolPtr(true),
    ListParams: sentinelone.ListParams{
        Limit:  sentinelone.IntPtr(50),
        SortBy: sentinelone.StringPtr("fullName"),
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%d users\n", pag.TotalItems)
```

### Get a single user

```go
user, err := client.Users.Get(ctx, "225494730938493801")
if err != nil {
    log.Fatal(err)
}
if user.Email != nil {
    fmt.Println(*user.Email, user.Scope)
}
```

### Create a user

```go
user, err := client.Users.Create(ctx, sentinelone.CreateUserRequest{
    Data: sentinelone.CreateUserData{
        Email:    "alice@example.com",
        FullName: "Alice Example",
        Scope:    "site",
        ScopeRoles: []sentinelone.UserScopeRole{
            {
                ID:          "225494730938493805", // site ID
                Name:        "Production East",
                AccountName: "Acme Corp",
                RoleID:      "225494730938493900", // role ID
            },
        },
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Created user:", user.ID)
```

### Update a user

```go
_, err = client.Users.Update(ctx, user.ID, sentinelone.UpdateUserRequest{
    Data: sentinelone.UpdateUserData{
        Scope:    "site",
        FullName: "Alice A. Example",
        ScopeRoles: []sentinelone.UserScopeRole{
            {ID: "225494730938493805", Name: "Production East", AccountName: "Acme Corp", RoleID: "225494730938493901"},
        },
    },
})
```

### Delete a user

```go
// Single user
err = client.Users.Delete(ctx, user.ID)

// Bulk delete by filter
err = client.Users.BulkDelete(ctx, sentinelone.BulkUsersActionRequest{
    Filter: sentinelone.BulkUsersFilter{
        IDs: []string{"225494730938493801", "225494730938493802"},
    },
})
```

### Two-factor authentication

```go
// Enable 2FA requirement for a user
err = client.Users.Enable2FA(ctx, sentinelone.UserIDRequest{
    Data: sentinelone.UserIDData{UserID: user.ID},
})

// Enroll the user (returns TOTP secret + QR code URL)
enroll, err := client.Users.Enroll2FA(ctx, sentinelone.UserIDsRequest{
    Data: sentinelone.UserIDsData{UserIDs: []string{user.ID}},
})
fmt.Println("TOTP secret:", enroll.Secret)

// Reset a user's 2FA device (e.g. lost phone)
err = client.Users.Reset2FA(ctx, sentinelone.ResetTFARequest{
    Data: sentinelone.ResetTFAData{UserID: user.ID},
})

// Disable 2FA requirement entirely
err = client.Users.Disable2FA(ctx, sentinelone.UserIDRequest{
    Data: sentinelone.UserIDData{UserID: user.ID},
})
```

### Password management

```go
// Force a password reset on next login for a set of users
err = client.Users.ForceResetPasswordOnLogin(ctx, sentinelone.ForceResetPasswordRequest{
    Filter: sentinelone.BulkUsersFilter{
        IDs: []string{user.ID},
    },
})

// Send a password reset email
err = client.Users.SendResetPasswordEmail(ctx, sentinelone.SendResetPasswordRequest{
    Filter: sentinelone.BulkUsersFilter{
        Email: "alice@example.com",
    },
})

// Change password (for the authenticated user)
err = client.Users.ChangePassword(ctx, sentinelone.ChangePasswordRequest{
    Data: sentinelone.ChangePasswordData{
        CurrentPassword: "old-password",
        NewPassword:     "new-password",
    },
})
```

### API token management

```go
// Check token metadata for a user
detail, err := client.Users.GetAPITokenDetails(ctx, user.ID)
fmt.Println("Token expires:", detail.ExpiresAt)

// Generate a new token for the authenticated user
tokenResp, err := client.Users.GenerateAPIToken(ctx, sentinelone.GenerateAPITokenRequest{})
fmt.Println("New API token:", tokenResp.Token)

// Revoke another user's token
err = client.Users.RevokeAPIToken(ctx, sentinelone.UserIDRequest{
    Data: sentinelone.UserIDData{UserID: user.ID},
})
```

### Authentication flows

```go
// Username + password login
loginResp, err := client.Users.Login(ctx, sentinelone.LoginRequest{
    Username: "admin@example.com",
    Password: "s3cur3P@ss",
})
if err != nil {
    log.Fatal(err)
}

// If 2FA is required, loginResp.Status == "2fa_required"
if loginResp.Status == "2fa_required" {
    loginResp2, err := client.Users.LoginContinue(ctx, sentinelone.LoginContinueRequest{
        Data: sentinelone.LoginContinueData{
            Token:  loginResp.Token,
            Code:   "123456", // TOTP code from authenticator app
            Method: loginResp.TwoFAMethod,
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    // Use loginResp2.Token for subsequent calls
    client = sentinelone.NewClient(baseURL, loginResp2.Token)
} else {
    client = sentinelone.NewClient(baseURL, loginResp.Token)
}

// Log out
err = client.Users.Logout(ctx)
```

---

## Complete example

The following program lists all active accounts and prints the sites and user count for each:

```go
package main

import (
    "context"
    "fmt"
    "log"

    s1 "github.com/s1buildpartners/sentinelone-go-sdk"
)

func main() {
    ctx := context.Background()

    client := s1.NewClient(
        "https://your-tenant.sentinelone.net",
        "your-api-token",
    )

    // Page through all active accounts
    var cursor *string
    for {
        accounts, pag, err := client.Accounts.List(ctx, &s1.ListAccountsParams{
            State: "active",
            ListParams: s1.ListParams{
                Limit:  s1.IntPtr(100),
                Cursor: cursor,
            },
        })
        if err != nil {
            log.Fatal(err)
        }

        for _, account := range accounts {
            fmt.Printf("Account: %s (ID: %s)\n", account.Name, account.ID)
            fmt.Printf("  Active agents : %d\n", account.ActiveAgents)
            fmt.Printf("  Sites         : %d\n", account.NumberOfSites)

            // List the first page of sites for this account
            siteResp, _, err := client.Sites.List(ctx, &s1.ListSitesParams{
                AccountID: account.ID,
                ListParams: s1.ListParams{Limit: s1.IntPtr(10)},
            })
            if err != nil {
                log.Printf("  warning: could not list sites: %v", err)
                continue
            }
            for _, site := range siteResp.Sites {
                fmt.Printf("    - %s (%s)\n", site.Name, site.State)
            }
        }

        if pag == nil || pag.NextCursor == nil {
            break
        }
        cursor = pag.NextCursor
    }
}
```

---

## API reference

Full GoDoc is available locally:

```bash
go doc github.com/s1buildpartners/sentinelone-go-sdk
```

The [official SentinelOne API documentation](https://your-tenant.sentinelone.net/api-doc/overview) lists valid enum values, permission requirements, and request/response field descriptions for each endpoint.
