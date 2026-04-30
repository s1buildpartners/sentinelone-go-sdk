package sentinelone

// UpdateSitesModulesRequest is the request body for
// PUT /web/api/v2.1/licenses/update-sites-modules.
type UpdateSitesModulesRequest struct {
	Data   UpdateSitesModulesData   `json:"data"`
	Filter UpdateSitesModulesFilter `json:"filter"`
}

// UpdateSitesModulesData specifies which operation to perform and which modules
// to add or remove.
type UpdateSitesModulesData struct {
	// Operation must be "add" or "remove".
	Operation string `json:"operation"`
	// Modules is the list of module names to add or remove.
	// Pass nil to affect the default module set for the operation.
	Modules []LicenseModuleItem `json:"modules,omitempty"`
}

// UpdateSitesModulesFilter restricts which sites are modified by
// [LicensesClient.UpdateSitesModules].
//
// At least one of SiteIDs or AccountIDs must be provided; the API
// rejects requests that do not target at least one site.
//
//   - SiteIDs / AccountIDs: select sites directly or by owning account.
//   - Query: full-text search on name, account name, and description.
//   - Name: exact site name match.
//   - State: "active", "expired", or "deleted".
//   - SiteType: "Trial" or "Paid".
type UpdateSitesModulesFilter struct {
	SiteIDs    []string `json:"siteIds,omitempty"`
	AccountIDs []string `json:"accountIds,omitempty"`
	Query      string   `json:"query,omitempty"`
	Name       string   `json:"name,omitempty"`
	State      string   `json:"state,omitempty"`
	SiteType   string   `json:"siteType,omitempty"`
}

// ---- License configuration request types ----

// LicensesInput is the license configuration block sent when creating or
// updating an account or site.  Set it on [CreateAccountData.Licenses],
// [UpdateAccountData.Licenses], [CreateSiteData.Licenses], or
// [UpdateSiteData.Licenses].
//
// Use either legacy SKUs or LicensesInput.Bundles, not both, in the same
// request.  Each bundle must include at least one surface entry.
type LicensesInput struct {
	Bundles []LicenseBundleInput `json:"bundles,omitempty"`
}

// LicenseBundleInput describes a single SKU and its configuration within a
// [LicensesInput] request.
//
// Name is required; use one of the Bundle* constants (e.g. [BundleComplete]).
// Surfaces is required; include at least one [LicenseSurfaceInput].
// Modules and Settings are optional per-bundle add-ons and platform overrides.
//
// MajorVersion is almost always omitted; the management console selects the
// latest available version automatically when it is absent.
type LicenseBundleInput struct {
	Name         string                `json:"name"`
	MajorVersion *int                  `json:"majorVersion,omitempty"`
	Surfaces     []LicenseSurfaceInput `json:"surfaces"`
	Modules      []LicenseModuleItem   `json:"modules,omitempty"`
	Settings     []LicenseSettingInput `json:"settings,omitempty"`
}

// LicenseSurfaceInput describes the entitlement count for a single deployment
// surface within a [LicenseBundleInput].
//
// Name must be one of the Surface* constants (e.g. [SurfaceTotalAgents]).
// Count is the numeric entitlement; use [SurfaceUnlimitedCount] (-1) for
// unlimited, otherwise set the specific numeric limit.
type LicenseSurfaceInput struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// LicenseModuleItem identifies an add-on module by name.
// Use one of the Module* constants (e.g. [ModuleSTAR]) for the Name field.
// LicenseModuleItem is used both inside [LicenseBundleInput.Modules] and in
// [UpdateSitesModulesData.Modules].
type LicenseModuleItem struct {
	Name string `json:"name"`
}

// LicenseSettingInput is a platform-level setting linked to a [LicenseBundleInput].
// GroupName must be one of the SettingGroup* constants;
// Setting must be one of the corresponding Setting* constants.
type LicenseSettingInput struct {
	GroupName string `json:"groupName"`
	Setting   string `json:"setting"`
}

// SurfaceUnlimitedCount is the sentinel value for [LicenseSurfaceInput.Count]
// that grants unlimited entitlement for that surface.
const SurfaceUnlimitedCount = -1

// ---- Bundle (SKU) name constants ----

// Endpoint base SKUs.
const (
	BundleCore     = "core"
	BundleControl  = "control"
	BundleComplete = "complete"
)

// Cloud Workload Security SKUs — per deployment surface.
const (
	BundleCWSServersControl               = "cloud_workload_security_servers_control"
	BundleCWSServersComplete              = "cloud_workload_security_servers_complete"
	BundleCWSContainersControl            = "cloud_workload_security_containers_control"
	BundleCWSContainersComplete           = "cloud_workload_security_containers_complete"
	BundleCWSServerlessContainersControl  = "cloud_workload_security_serverless_containers_control"
	BundleCWSServerlessContainersComplete = "cloud_workload_security_serverless_containers_complete"
)

// Data and logging SKUs.
const (
	BundleSingularityDataLake = "singularity_data_lake"
	BundleLogAnalytics        = "log_analytics"
)

// Identity SKUs.
const (
	BundleIdentityThreatProtection = "identity_threat_protection"
	BundleUnifiedIdentity          = "unified_identity"
	BundleSingularityIdentity      = "singularity_identity"
	BundleRangerADProtect          = "ranger_ad_protect"
)

// Other base service SKUs.
const (
	BundleMobile                = "mobile"
	BundleThreatDetectionS3     = "threat_detection_s3"
	BundleThreatDetectionNetApp = "threat_detection_netapp"
)

// ---- Surface name constants ----

const (
	// SurfaceTotalAgents is the surface for endpoint and CWS SKUs.
	SurfaceTotalAgents = "Total Agents"

	// SurfaceTotalEndpoints is the surface for identity and datastore SKUs.
	SurfaceTotalEndpoints = "Total Endpoints"

	// SurfaceTotalUsers is the surface for user-based identity SKUs.
	SurfaceTotalUsers = "Total Users"

	// SurfaceAverageGBPerDay is the data-volume surface for SDL-style SKUs.
	SurfaceAverageGBPerDay = "Average GB/Day"

	// SurfaceLongRangeQueryCredits is the credit-based query surface used
	// alongside singularity_data_lake and related SKUs.
	SurfaceLongRangeQueryCredits = "Long-Range Query Credits"
)

// ---- Add-on module (modules[].name) constants ----

// Classic endpoint add-ons.
const (
	ModuleRSO                  = "rso"
	ModuleSTAR                 = "star"
	ModuleBinaryVaultBenign    = "binary_vault_benign"
	ModuleBinaryVaultMalicious = "binary_vault_malicious"
	ModuleRanger               = "ranger"
	ModuleRogues               = "rogues"
)

// MDR / vigilance / alerts add-ons.
const (
	ModuleVigilance              = "vigilance"
	ModuleWatchtower             = "watchtower"
	ModuleSingularityMDR         = "singularity_mdr"
	ModuleSkylightAlerts         = "skylight_alerts"
	ModuleWayfinderMDREssentials = "wayfinder_mdr_essentials"
	ModuleWayfinderMDRElite      = "wayfinder_mdr_elite"
	ModuleWayfinderThreatHunting = "wayfinder_threat_hunting"
)

// Automation / XDR / extended retention add-ons.
const (
	ModuleHyperautomation   = "hyperautomation"
	ModuleXDR1095DRetention = "xdr_1095d_retention"
	ModuleXDR1460DRetention = "xdr_1460d_retention"
)

// Setting group names for use in [LicenseSettingInput.GroupName].
const (
	SettingGroupMaliciousDataRetention  = "malicious_data_retention"
	SettingGroupDVRetention             = "dv_retention"
	SettingGroupRemoteShellAvailability = "remote_shell_availability"
	SettingGroupMarketplaceAccessStatus = "marketplace_access_status"
	SettingGroupAccountLevelRanger      = "account_level_ranger"
)

// ---- Setting value constants ----

// Values for [SettingGroupMaliciousDataRetention].
const (
	SettingMaliciousRetention365Days = "365 Days"
)

// Values for [SettingGroupDVRetention].
const (
	SettingDVRetention4Days   = "4 Days"
	SettingDVRetention30Days  = "30 Days"
	SettingDVRetention90Days  = "90 Days"
	SettingDVRetention180Days = "180 Days"
	SettingDVRetention365Days = "365 Days"
)

// Values for [SettingGroupRemoteShellAvailability].
const (
	SettingRemoteShellEnabled  = "Enabled"
	SettingRemoteShellDisabled = "Disabled"
)

// Values for [SettingGroupMarketplaceAccessStatus].
const (
	SettingMarketplaceAvailable = "Available"
	SettingMarketplaceNoAccess  = "No Access"
)

// Values for [SettingGroupAccountLevelRanger].
const (
	SettingRangerLevelAccount = "Account"
	SettingRangerLevelSite    = "Site"
)
