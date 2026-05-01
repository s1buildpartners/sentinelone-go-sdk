# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1] - 2026-04-30

### Added

- CodeQL analysis workflow (`.github/workflows/codeql.yml`) — runs on push and
  pull requests to `main` and on a weekly schedule; covers Go source with the
  CodeQL security and quality query suite.

### Fixed

- `LicenseBundleDataIngest` constant: corrected API value from `"data_ingest"`
  to `"singularity_data_lake"`.
- `NewMarketplaceAccessSettingInput`: corrected the enabled-state API value from
  `"Enabled"` to `"Available"`.
- `NewXDRDataRetentionSettingInput`: corrected boundary conditions in the
  day-to-setting mapping.  The previous `>=` comparisons caused values at exact
  thresholds (e.g. 30, 90, 180) to map to the wrong tier; the corrected `>`
  comparisons match the documented "nearest valid maximum" semantics (e.g. 30
  days → `"30 Days"`, 31 days → `"90 Days"`).
- Release workflow: fixed changelog-section extraction step that was using the
  wrong version variable, and added a step that waits for the CodeQL workflow to
  complete before creating the GitHub release.

## [0.2.0] - 2026-04-30

### Added

- Named constants for account and site type fields:
  - `AccountTypeTrial`, `AccountTypePaid` — valid values for `AccountType` on
    `CreateAccountData`, `UpdateAccountData`, and `ListAccountsParams`
  - `SiteTypeTrial`, `SiteTypePaid` — valid values for `SiteType` on
    `CreateSiteData`, `UpdateSiteData`, `ListSitesParams`, and
    `UpdateSitesModulesFilter`
- Named constants for entity state filter fields (`StateActive`, `StateExpired`,
  `StateDeleted`) — valid values for `State` on `ListAccountsParams`,
  `ListSitesParams`, and `UpdateSitesModulesFilter`
- Named constants for license setting values:
  - `LicenseSettingNetworkDiscoveryConsolidationLevelAccount`,
    `LicenseSettingNetworkDiscoveryConsolidationLevelSite` — valid values for
    `NewNetworkDiscoveryConsolidationLevelSettingInput`
  - `LicenseSettingIdentitySecurityPostureModeFull`,
    `LicenseSettingIdentitySecurityPostureModeLite` — valid values for
    `NewIdentitySecurityPostureModeSettingInput`
- Additional bundle, module, and setting constants (expanding the set from
  `[0.1.0]`):
  - `LicenseBundle*` — 22 SKU names covering Endpoint Security (Core, Complete,
    Control), Cloud Workload Security (servers, serverless containers, containers
    — Control and Complete tiers each), Cloud-Native Security (Pro, Foundations),
    Data Ingest, Log Analytics, Mobile Security, Hyperautomation, Identity Threat
    Detection, Identity Security Posture Management, Identity Security for IDP,
    Unified Identity, Identity Detection and Response, and Threat Detection for
    NetApp and Data Stores
  - `LicenseModule*` — 30 add-on names covering XDR data retention tiers (30 d –
    5 y), Cloud Funnel, Vigilance MDR, Binary Vault, Network Discovery, Singularity
    MDR, Threat Intel, Purple AI (Foundations, SOC Analyst), WatchTower,
    Vulnerability Management, Unprotected Endpoint Discovery, Remote Script
    Orchestration, Wayfinder (Elite, Essentials, Threat Hunting), and Remote Ops
    Forensics
  - `LicenseSurface*` — 7 surface names: `TotalAgents`, `TotalEndpoints`,
    `TotalUsers`, `AvgGBDay`, `Workloads`, `LongRangeQueryCredits`,
    `ActionPacks`; plus `LicenseSurfaceUnlimitedCount = -1`
  - `LicenseSetting*` — 6 setting group names and 4 setting value constants
- Helper constructor functions that build correctly-typed `LicenseBundleInput`,
  `LicenseModuleItem`, and `LicenseSettingInput` values without manual struct
  construction:
  - **Bundle constructors** (one per SKU, named `New<SKU>BundleInput`) — accept
    the relevant entitlement count(s) and return a `LicenseBundleInput` with the
    correct name and surfaces pre-populated
  - **Module constructors** (`New<Module>ModuleItem`) — return a
    `LicenseModuleItem` for each supported add-on
  - **Setting constructors** — `NewXDRDataRetentionSettingInput(days int)`,
    `NewEDRDataRetentionSettingInput()`, `NewRemoteShellSettingInput(enabled bool)`,
    `NewMarketplaceAccessSettingInput(enabled bool)`,
    `NewNetworkDiscoveryConsolidationLevelSettingInput(level string)`,
    `NewIdentitySecurityPostureModeSettingInput(mode string)`

- `WithLogger(l *slog.Logger) ClientOption` — injects a structured logger into
  the client.  By default all output is discarded (`slog.DiscardHandler`).  Pass
  your application's logger to surface request traces, retry warnings, and API
  error details without touching the client internals.
- Structured logging throughout the central `do()` request loop:
  - `Debug` — rate-limit token wait, each outbound request (method, path,
    attempt number), and successful response completion (status code)
  - `Warn` — rate-limit wait interrupted by context cancellation, 429 backoff
    (method, path, attempt, `retry_after` duration), and any non-2xx API
    response (method, path, status code)
  - `Error` — request body marshal failure, request build failure, HTTP
    transport error, response body read error, and response decode failure

### Removed

- Inaccurate `// Required permission:` doc-comment lines from all API methods
  in `accounts.go`, `agents.go`, `licenses.go`, `rbac.go`, `sites.go`, and
  `users.go`.  The permission names were not verified against the live API and
  were therefore misleading.  Refer to the SentinelOne API Hub in your console
  for authoritative permission information.

### Changed

- License input types overhauled (breaking changes relative to `[0.1.0]`):
  - All license enum constants renamed with the `License` prefix (e.g.
    `BundleComplete` → `LicenseBundleEndpointSecurityComplete`, `ModuleRanger` →
    `LicenseModuleNetworkDiscovery`, `SettingGroupDVRetention` →
    `LicenseSettingXDRDataRetention`)
  - `LicenseBundleInput` simplified: the `Modules` and `Settings` fields have
    moved up to `LicensesInput`; `LicenseBundleInput` now carries only `Name`,
    `MajorVersion`, and `Surfaces`
  - `doc.go` and `README.md` examples updated throughout to use the new constant
    names and helper constructor functions instead of manual struct literals
- `README.md` — added a Table of Contents, a Developer Notes section listing
  the toolchain used by maintainers, and a Questions/Issues/Feature Requests
  section.  The API-documentation reference now points to **Help > API Hub** in
  the SentinelOne console rather than a placeholder URL.

## [0.1.0] - 2026-04-29

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

[Unreleased]: https://github.com/s1buildpartners/sentinelone-go-sdk/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/s1buildpartners/sentinelone-go-sdk/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/s1buildpartners/sentinelone-go-sdk/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/s1buildpartners/sentinelone-go-sdk/releases/tag/v0.1.0
