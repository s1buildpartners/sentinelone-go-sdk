# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### Package structure

- `types/` subpackage (`package types`) containing all domain and response types,
  following the AWS SDK Go v2 pattern. Input types (request bodies, filter params)
  live in the root package; model types (API response structs) live in `types/`.
  - `types/accounts.go` — `Account`, `AccountSKU`, `UninstallPasswordMetadata`,
    `UninstallPassword`
  - `types/common.go` — `Policy`, `PolicyEngines`, `Licenses`, `IRFields` and
    related embedded types
  - `types/errors.go` — `APIError`, `ResponseError`, `Pagination`
  - `types/rbac.go` — `Role`, `RolePage`, `RolePermissionEntry`,
    `RoleWithPermissions`
  - `types/sites.go` — `Site`, `SitesResponse`, `SiteToken`, `SiteKey`,
    `LocalAuthorization`
  - `types/users.go` — `User`, `UserAPIToken`, `UserScopeRole`, `UserSiteRole`,
    and all authentication response types (`LoginResponse`,
    `LoginContinueResponse`, `APITokenResponse`, `APITokenDetail`,
    `EnrollTFAResponse`, `IFrameTokenResponse`, `ElevateSessionResponse`,
    `RequestAppResponse`, `SetPasswordResponse`)

#### Client

- `Client` with `Accounts`, `Sites`, `RBAC`, and `Users` sub-client fields
- `NewClient(baseURL, apiToken, ...ClientOption)` constructor
- `WithTimeout(d)` option to override the default 30-second per-request timeout
- `WithHTTPClient(hc)` option to provide a custom `*http.Client`
- `WithRateLimiting(bool)` option to enable or disable the built-in rate limiter
  (enabled by default)
- `WithMaxRetries(n)` option to set the maximum number of automatic 429 retries
  (default: 3)
- `AsResponseError(err)` helper to unwrap `*types.ResponseError` without
  importing the `errors` package directly
- `BoolPtr`, `IntPtr`, `StringPtr` pointer-constructor helpers for optional
  request fields
- `ListParams` struct with cursor-based pagination fields (`Cursor`, `Limit`,
  `Skip`, `SortBy`, `SortOrder`, `CountOnly`, `SkipCount`) embedded by all
  endpoint-specific params types

#### Rate limiting

- Per-path proactive token-bucket rate limiting covering all documented
  SentinelOne MGMT API v2.1 per-API-token rate limits (43 path prefixes),
  using longest-prefix matching so sub-paths receive the correct independent
  limit (e.g. `/threats/<id>/notes` at 100 req/s vs `/threats` at 10 req/s)
- Automatic reactive retry on HTTP 429 Too Many Requests: honours the
  `Retry-After` response header and falls back to a 5-second wait when the
  header is absent; context cancellation aborts the backoff immediately

#### API groups

- `AccountsClient` — `List`, `Get`, `Create`, `Update`, `ExpireNow`,
  `Reactivate`, `GetPolicy`, `UpdatePolicy`, `RevertPolicy`,
  `GetUninstallPasswordMetadata`, `GenerateUninstallPassword`,
  `RevokeUninstallPassword`
- `SitesClient` — `List`, `Get`, `Create`, `Update`, `Delete`, `BulkUpdate`,
  `Duplicate`, `Reactivate`, `GetPolicy`, `UpdatePolicy`, `RevertPolicy`,
  `GetToken`, `RegenerateKey`, `GetLocalAuthorization`,
  `UpdateLocalAuthorization`
- `RBACClient` — `List`, `GetTemplate`, `Get`, `Create`, `Update`, `Delete`
- `UsersClient` — `List`, `Get`, `Create`, `Update`, `Delete`, `BulkDelete`,
  `BulkEnable`, `BulkDisable`; authentication: `Login`, `LoginContinue`,
  `LoginByToken`, `Logout`, `SSOReAuth`, `ElevateSession`, `SetPassword`;
  API tokens: `GenerateAPIToken`, `GetAPITokenDetails`, `RevokeAPIToken`;
  2FA: `Enroll2FA`, `Enable2FA`, `Disable2FA`, `Reset2FA`, `DeleteTFA`;
  passwords: `ChangePassword`, `SendResetPasswordEmail`,
  `ForceResetPasswordOnLogin`; misc: `SendVerificationEmail`,
  `GenerateIFrameToken`, `RequestApp`, `EnableApp`,
  `OnboardingValidateToken`, `OnboardingVerifyToken`

#### CI/CD

- `.github/workflows/main.yml` — CI workflow triggered on push and pull
  requests to `main`: runs markdownlint, golangci-lint v2, and
  `go test -race` with coverage reporting
- `.github/workflows/release.yml` — release workflow triggered on
  `v*.*.*` tags: extracts the matching section from this changelog and
  creates a GitHub release with those notes; pre-release detection via
  tag suffix (e.g. `-beta.1`)

[Unreleased]: https://github.com/s1buildpartners/sentinelone-go-sdk/commits/main/
