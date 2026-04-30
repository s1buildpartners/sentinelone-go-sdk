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
  - `types/agents.go` — `Agent`, `AgentNetworkInterface`, `AgentLocation`,
    `VssVolume`, `DiskMetrics`
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

- `Client` with `Accounts`, `Sites`, `RBAC`, `Users`, `Agents`, and `Licenses`
  sub-client fields
- `NewClient(baseURL, apiToken, ...ClientOption)` constructor
- `NewClientFromEnv()` constructor — reads `SENTINELONE_URL` and
  `SENTINELONE_TOKEN` from the environment; returns an error if either is absent
- `NewClientFromConfig(...LoadOption)` constructor — reads credentials from an
  INI-format file under a named profile; supports `WithProfile` and
  `WithConfigFile` options
- `NewDefaultClient(...LoadOption)` constructor — layered credential lookup:
  environment variables take priority, falling back to the credentials file
- `WithProfile(name)` load option — selects a named profile in the credentials
  file (falls back to `SENTINELONE_PROFILE`, then `"default"`)
- `WithConfigFile(path)` load option — overrides the platform-default credentials
  file path and the `SENTINELONE_CONFIG` environment variable
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
- INI-format credentials file with per-profile `url` and `token` keys;
  platform-default paths: `$XDG_CONFIG_HOME/sentinelone/credentials` (Linux/BSD),
  `~/Library/Application Support/sentinelone/credentials` (macOS),
  `%AppData%\SentinelOne\credentials` (Windows)

#### Rate limiting

- Per-path proactive token-bucket rate limiting covering all documented
  SentinelOne MGMT API v2.1 per-API-token rate limits (43 path prefixes),
  using longest-prefix matching so sub-paths receive the correct independent
  limit (e.g. `/threats/<id>/notes` at 100 req/s vs `/threats` at 10 req/s)
- Automatic reactive retry on HTTP 429 Too Many Requests: honours the
  `Retry-After` response header and falls back to a 5-second wait when the
  header is absent; context cancellation aborts the backoff immediately

#### API groups

- `AccountsClient` — `List`, `Get`, `Create`, `Update`, `UpdateLicenses`,
  `ExpireNow`, `Reactivate`, `GetPolicy`, `UpdatePolicy`, `RevertPolicy`,
  `GetUninstallPasswordMetadata`, `GenerateUninstallPassword`,
  `RevokeUninstallPassword`
- `SitesClient` — `List`, `Get`, `Create`, `Update`, `UpdateLicenses`, `Delete`,
  `BulkUpdate`, `Duplicate`, `Reactivate`, `GetPolicy`, `UpdatePolicy`,
  `RevertPolicy`, `GetToken`, `RegenerateKey`, `GetLocalAuthorization`,
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
- `AgentsClient` — `List`, `Count`; `ListAgentsParams` with 50+ filter fields
  covering identity, OS type, version, network status, scan state, mitigation
  mode, threat count, operational state, and date-range filters
- `LicensesClient` — `UpdateSitesModules` (bulk add/remove add-on modules across
  sites matching a filter)

#### License configuration

- `LicensesInput` — top-level license block for create/update requests
- `LicenseBundleInput` — one SKU with surfaces, modules, and per-bundle settings
- `LicenseSurfaceInput` — entitlement count for a single deployment surface
  (`Count: -1` / `SurfaceUnlimitedCount` for unlimited)
- `LicenseModuleItem` — add-on module referenced by name (shared by bundle
  definition and `UpdateSitesModulesData`)
- `LicenseSettingInput` — platform setting associated with a bundle
- `Licenses *LicensesInput` field added to `CreateAccountData`,
  `UpdateAccountData`, `CreateSiteData`, and `UpdateSiteData`
- Named constants for all known enum values:
  - `Bundle*` — 17 SKU names (`BundleCore`, `BundleComplete`,
    `BundleCWSServersControl`, `BundleSingularityDataLake`, …)
  - `Surface*` — 5 surface names (`SurfaceTotalAgents`, `SurfaceTotalUsers`,
    `SurfaceAverageGBPerDay`, …) plus `SurfaceUnlimitedCount = -1`
  - `Module*` — 17 add-on names (`ModuleSTAR`, `ModuleRanger`,
    `ModuleVigilance`, `ModuleXDR1095DRetention`, …)
  - `SettingGroup*` — 5 setting group names (`SettingGroupDVRetention`,
    `SettingGroupAccountLevelRanger`, …)
  - `Setting*` — setting values per group (`SettingDVRetention90Days`,
    `SettingRangerLevelAccount`, …)

#### CI/CD

- `.github/workflows/main.yml` — CI workflow triggered on push and pull
  requests to `main`: runs markdownlint, golangci-lint v2, and
  `go test -race` with coverage reporting
- `.github/workflows/release.yml` — release workflow triggered on
  `v*.*.*` tags: extracts the matching section from this changelog and
  creates a GitHub release with those notes; pre-release detection via
  tag suffix (e.g. `-beta.1`)
- `.github/workflows/main.yml` — added `govulncheck` job that runs
  `golang.org/x/vuln/cmd/govulncheck` against the module on every push and
  pull request

### Changed

- `NewClientFromProfile` renamed to `NewDefaultClient` to better reflect its
  layered credential lookup behaviour (environment variables → credentials file)

### Fixed

- `govulncheck` CI job: corrected workflow step configuration that prevented the
  job from running successfully

[Unreleased]: https://github.com/s1buildpartners/sentinelone-go-sdk/commits/main/
