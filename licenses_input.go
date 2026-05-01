package sentinelone

import "strings"

const (
	// LicenseSurfaceUnlimitedCount is the sentinel value for an unlimited entitlement count in a
	// [LicenseSurfaceInput].
	LicenseSurfaceUnlimitedCount = -1
)

// License surface strings must match exactly the API's expected values; use the constants provided in this package
// (e.g. [LicenseSurfaceTotalAgents]) rather than hardcoding strings in your code.  The API is case-sensitive and
// rejects requests with unrecognized surface names, so these constants must be used as-is without modification
// (e.g. no title-casing).
const (
	LicenseSurfaceTotalAgents           = "Total Agents"
	LicenseSurfaceTotalEndpoints        = "Total Endpoints"
	LicenseSurfaceTotalUsers            = "Total Users"
	LicenseSurfaceAvgGBDay              = "Average GB per Day"
	LicenseSurfaceWorkloads             = "Workloads"
	LicenseSurfaceLongRangeQueryCredits = "Long Range Query Credits"
	LicenseSurfaceActionPacks           = "Action Packs"
)

// License bundle names must match exactly the API's expected values; use the constants provided in this package
// (e.g. [LicenseBundleEndpointSecurityComplete]) rather than hardcoding strings in your code.  The API is
// case-sensitive and rejects requests with unrecognized bundle names, so these constants must be used as-is without
// modification (e.g. no title-casing).
//
// Use the New*BundleInput functions (e.g. [NewEndpointSecurityCompleteBundleInput]) to create bundles with the
// correct name and surfaces.
const (
	LicenseBundleThreatDetectionNetApp              = "threat_detection_netapp"
	LicenseBundleThreatDetectionDataStores          = "threat_detection_s3"
	LicenseBundleMobileSecurity                     = "mobile_security"
	LicenseBundleLogAnalytics                       = "log_analytics"
	LicenseBundleIdentityThreatProtection           = "identity_threat_protection"
	LicenseBundleIdentitySecurityPostureManagement  = "ranger_ad"
	LicenseBundleIdentitySecurityForIDP             = "ranger_ad_protect"
	LicenseBundleUnifiedIdentity                    = "unified_identity"
	LicenseBundleIdentityDetectionResponse          = "singularity_identity"
	LicenseBundleHyperautomation                    = "hyperautomation"
	LicenseBundleEndpointSecurityCore               = "core"
	LicenseBundleEndpointSecurityComplete           = "complete"
	LicenseBundleEndpointSecurityControl            = "control"
	LicenseBundleDataIngest                         = "singularity_data_lake"
	LicenseBundleCWSForServersControl               = "cloud_workload_security_servers_control"
	LicenseBundleCWSForServersComplete              = "cloud_workload_security_servers_complete"
	LicenseBundleCWSForServerlessContainersControl  = "cloud_workload_security_serverless_containers_control"
	LicenseBundleCWSForServerlessContainersComplete = "cloud_workload_security_serverless_containers_complete"
	LicenseBundleCWSForContainersControl            = "cloud_workload_security_containers_control"
	LicenseBundleCWSForContainersComplete           = "cloud_workload_security_containers_complete"
	LicenseBundleCNSPro                             = "cloud_native_security_pro"
	LicenseBundleCNSFoundations                     = "cloud_native_security_foundations"
)

// License module names must match exactly the API's expected values; use the constants provided in this package
// (e.g. [LicenseModuleDataIngest30d]) rather than hardcoding strings in your code.  The API is case-sensitive and
// rejects requests with unrecognized module names, so these constants must be used as-is without modification
// (e.g. no title-casing).
//
// Use the New*ModuleInput functions (e.g. [NewDataIngestModuleInput]) to create modules with the correct name.
const (
	LicenseModuleDataIngest30d                     = "xdr_30d_retention"
	LicenseModuleDataIngest90d                     = "xdr_90d_retention"
	LicenseModuleDataIngest180d                    = "xdr_180d_retention"
	LicenseModuleDataIngest365d                    = "xdr_365d_retention"
	LicenseModuleDataIngestLongRangeRetention1y    = "xdr_extended_365d_retention"
	LicenseModuleDataIngestLongRangeRetention2y    = "xdr_extended_730d_retention"
	LicenseModuleDataIngestLongRangeRetention3y    = "xdr_extended_1095d_retention"
	LicenseModuleDataIngestLongRangeRetention4y    = "xdr_extended_1460d_retention"
	LicenseModuleDataIngestLongRangeRetention5y    = "xdr_extended_1825d_retention"
	LicenseModuleCloudFunnel                       = "cloud_funnel"
	LicenseModuleVigilanceMDR                      = "vigilance"
	LicenseModuleBinaryVaultBenignFiles            = "binary_vault_benign"
	LicenseModuleDataIngestLongRangeEndpointAndCWS = "xdr_edr_cws_retention"
	LicenseModuleNetworkDiscovery                  = "ranger"
	LicenseModuleSingularityMDR                    = "singularity_mdr"
	LicenseModuleThreatIntel                       = "threat_intel"
	LicenseModulePurpleAIFoundations               = "purple_ai"
	LicenseModulePurpleAISocAnalyst                = "purple_ai_soc_analyst"
	LicenseModuleWatchTower                        = "watchtower"
	LicenseModuleVulnerabilityManagement           = "vulnerability_management"
	LicenseModuleUnprotectedEndpointDiscovery      = "rogues"
	LicenseModuleRemoteScriptOrchestration         = "rso"
	LicenseModuleWayfinderElite                    = "wayfinder_mdr_elite"
	LicenseModuleWayfinderEssentials               = "wayfinder_mdr_essentials"
	LicenseModuleWayfinderThreatHunting            = "wayfinder_threat_hunting"
	LicenseModuleRemoteOpsForensics                = "remote_ops_forensics"
)

// License setting group and setting names must match exactly the API's expected values; use the constants provided
// in this package.  The API is case-sensitive and rejects requests with unrecognized setting group or setting names,
// so these constants must be used as-is without modification (e.g. no title-casing).
//
// Use the New*SettingInput functions (e.g. [NewXDRDataRetentionSettingInput]) to create settings with the correct
// group and setting names and values.
const (
	LicenseSettingXDRDataRetention                   = "dv_retention"
	LicenseSettingEDRDataRetention                   = "malicious_data_retention"
	LicenseSettingRemoteShell                        = "remote_shell_availability"
	LicenseSettingMarketplaceAccess                  = "marketplace_access_status"
	LicenseSettingNetworkDiscoveryConsolidationLevel = "account_level_ranger"
	LicenseSettingIdentitySecurityPostureMode        = "ranger_ad_mode"
)

// License setting value constants for [LicenseSettingNetworkDiscoveryConsolidationLevel] and
// [LicenseSettingIdentitySecurityPostureMode].
const (
	LicenseSettingNetworkDiscoveryConsolidationLevelAccount = "Account"
	LicenseSettingNetworkDiscoveryConsolidationLevelSite    = "Site"
	LicenseSettingIdentitySecurityPostureModeFull           = "Full"
	LicenseSettingIdentitySecurityPostureModeLite           = "Lite"
)

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
//   - State: [StateActive], [StateExpired], or [StateDeleted].
//   - SiteType: [SiteTypeTrial] or [SiteTypePaid].
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
	Bundles          []LicenseBundleInput  `json:"bundles,omitempty"`
	MakeSocDefaultUI *bool                 `json:"makeSocDefaultUi,omitempty"`
	Modules          []LicenseModuleItem   `json:"modules,omitempty"`
	Name             *string               `json:"name,omitempty"`
	Settings         []LicenseSettingInput `json:"settings,omitempty"`
}

// LicenseBundleInput describes a single SKU and its configuration within a
// [LicensesInput] request.
//
// Name is required; use one of the Bundle* variables (e.g. [BundleComplete]).
// Surfaces is required; include at least one [LicenseSurfaceInput].
// Modules and Settings are optional per-bundle add-ons and platform overrides.
//
// MajorVersion and MinorVersion are almost always omitted; the management
// console selects the latest available version automatically when absent.
type LicenseBundleInput struct {
	Name         string                `json:"name"`
	MajorVersion *int                  `json:"majorVersion,omitempty"`
	Surfaces     []LicenseSurfaceInput `json:"surfaces"`
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
// Use one of the Module* variables (e.g. [ModuleSTAR]) for the Name field.
// LicenseModuleItem is used both inside [LicenseBundleInput.Modules] and in
// [UpdateSitesModulesData.Modules].
type LicenseModuleItem struct {
	Name string `json:"name"`
}

// LicenseSettingInput is a platform-level setting linked to a [LicenseBundleInput].
// GroupName must be one of the SettingGroup* variables;
// Setting must be one of the corresponding Setting* constants.
type LicenseSettingInput struct {
	GroupName string `json:"groupName"`
	Setting   string `json:"setting"`
}

// NewThreatDetectionForNetAppBundleInput creates a [LicenseBundleInput] for the Threat Detection for NetApp bundle
// with the specified total endpoints.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited endpoints.
func NewThreatDetectionForNetAppBundleInput(totalEndpoints int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalEndpoints,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleThreatDetectionNetApp,
		Surfaces: surfaces,
	}
}

// NewThreatDetectionForDataStoresBundleInput creates a [LicenseBundleInput] for the Threat Detection for Data Stores
// bundle with the specified total endpoints.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited endpoints.
func NewThreatDetectionForDataStoresBundleInput(totalEndpoints int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalEndpoints,
			Count: totalEndpoints,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleThreatDetectionDataStores,
		Surfaces: surfaces,
	}
}

// NewMobileSecurityBundleInput creates a [LicenseBundleInput] for the Mobile Security bundle with the specified
// total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewMobileSecurityBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleMobileSecurity,
		Surfaces: surfaces,
	}
}

// NewLogAnalyticsBundleInput creates a [LicenseBundleInput] for the Log Analytics bundle with the specified average
// GB/day and long-range query credits.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited average GB/day or long-range query credits.
//
// Note that this bundle is deprecated and may not be available in all accounts; check with SentinelOne support
// before using it.
func NewLogAnalyticsBundleInput(avgGBPerDay, longRangeQueryCredits int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceAvgGBDay,
			Count: avgGBPerDay,
		},
		{
			Name:  LicenseSurfaceLongRangeQueryCredits,
			Count: longRangeQueryCredits,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleLogAnalytics,
		Surfaces: surfaces,
	}
}

// NewIdentityThreatDetectionBundleInput creates a [LicenseBundleInput] for the Identity Threat Detection bundle with
// the specified total endpoints.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited endpoints.
func NewIdentityThreatDetectionBundleInput(totalEndpoints int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalEndpoints,
			Count: totalEndpoints,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleIdentityThreatProtection,
		Surfaces: surfaces,
	}
}

// NewIdentitySecurityPostureManagementBundleInput creates a [LicenseBundleInput] for the Identity Security Posture
// Management bundle with the specified total endpoints.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited users.
func NewIdentitySecurityPostureManagementBundleInput(totalUsers int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalUsers,
			Count: totalUsers,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleIdentitySecurityPostureManagement,
		Surfaces: surfaces,
	}
}

// NewIdentitySecurityForIDPBundleInput creates a [LicenseBundleInput] for the Identity Security for IDP bundle with
// the specified total users.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited users.
func NewIdentitySecurityForIDPBundleInput(totalUsers int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalUsers,
			Count: totalUsers,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleIdentitySecurityForIDP,
		Surfaces: surfaces,
	}
}

// NewIdentitySecurityBundleInput creates a [LicenseBundleInput] for the Unified Identity bundle with the specified
// total users.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited users.
func NewIdentitySecurityBundleInput(totalUsers int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalUsers,
			Count: totalUsers,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleUnifiedIdentity,
		Surfaces: surfaces,
	}
}

// NewIdentityDetectionResponseBundleInput creates a [LicenseBundleInput] for the Identity Detection and Response
// bundle with the specified total endpoints.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited endpoints.
func NewIdentityDetectionResponseBundleInput(totalEndpoints int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalEndpoints,
			Count: totalEndpoints,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleIdentityDetectionResponse,
		Surfaces: surfaces,
	}
}

// NewHyperautomationBundleInput creates a [LicenseBundleInput] for the Hyperautomation bundle with the specified
// total action packs.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited action packs.
func NewHyperautomationBundleInput(totalActionPacks int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceActionPacks,
			Count: totalActionPacks,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleHyperautomation,
		Surfaces: surfaces,
	}
}

// NewEndpointSecurityCoreBundleInput creates a [LicenseBundleInput] for the Endpoint Security Core bundle with
// the specified total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewEndpointSecurityCoreBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleEndpointSecurityCore,
		Surfaces: surfaces,
	}
}

// NewEndpointSecurityCompleteBundleInput creates a [LicenseBundleInput] for the Endpoint Security Complete bundle
// with the specified total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewEndpointSecurityCompleteBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleEndpointSecurityComplete,
		Surfaces: surfaces,
	}
}

// NewEndpointSecurityControlBundleInput creates a [LicenseBundleInput] for the Endpoint Security Control bundle
// with the specified total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewEndpointSecurityControlBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleEndpointSecurityControl,
		Surfaces: surfaces,
	}
}

// NewDataIngestBundleInput creates a [LicenseBundleInput] for the Data Ingest bundle with the specified average
// GB/day and long-range query credits.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited average GB/day or long-range query credits.
func NewDataIngestBundleInput(avgGBPerDay, longRangeQueryCredits int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceAvgGBDay,
			Count: avgGBPerDay,
		},
		{
			Name:  LicenseSurfaceLongRangeQueryCredits,
			Count: longRangeQueryCredits,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleDataIngest,
		Surfaces: surfaces,
	}
}

// NewCWSForServersControlBundleInput creates a [LicenseBundleInput] for the Cloud Workload Security for
// Servers Control bundle with the specified total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewCWSForServersControlBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleCWSForServersControl,
		Surfaces: surfaces,
	}
}

// NewCWSForServersCompleteBundleInput creates a [LicenseBundleInput] for the Cloud Workload Security for
// Servers Complete bundle with the specified total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewCWSForServersCompleteBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleCWSForServersComplete,
		Surfaces: surfaces,
	}
}

// NewCWSForServerlessContainersControlBundleInput creates a [LicenseBundleInput] for the Cloud Workload Security for
// Serverless Containers Control bundle with the specified total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewCWSForServerlessContainersControlBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleCWSForServerlessContainersControl,
		Surfaces: surfaces,
	}
}

// NewCWSForServerlessContainersCompleteBundleInput creates a [LicenseBundleInput] for the Cloud Workload Security for
// Serverless Containers Complete bundle with the specified total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewCWSForServerlessContainersCompleteBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleCWSForServerlessContainersComplete,
		Surfaces: surfaces,
	}
}

// NewCWSForContainersControlBundleInput creates a [LicenseBundleInput] for the Cloud Workload Security for
// Containers Control bundle with the specified total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewCWSForContainersControlBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleCWSForContainersControl,
		Surfaces: surfaces,
	}
}

// NewCWSForContainersCompleteBundleInput creates a [LicenseBundleInput] for the Cloud Workload Security for
// Containers Complete bundle with the specified total agents.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited agents.
func NewCWSForContainersCompleteBundleInput(totalAgents int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceTotalAgents,
			Count: totalAgents,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleCWSForContainersComplete,
		Surfaces: surfaces,
	}
}

// NewCNSProBundleInput creates a [LicenseBundleInput] for the Cloud Native Security Pro bundle with the specified
// total workloads.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited workloads.
func NewCNSProBundleInput(totalWorkloads int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceWorkloads,
			Count: totalWorkloads,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleCNSPro,
		Surfaces: surfaces,
	}
}

// NewCNSFoundationsBundleInput creates a [LicenseBundleInput] for the Cloud Native Security Foundations bundle with
// the specified total workloads.
//
// Use [LicenseSurfaceUnlimitedCount] for unlimited workloads.
func NewCNSFoundationsBundleInput(totalWorkloads int) LicenseBundleInput {
	surfaces := []LicenseSurfaceInput{
		{
			Name:  LicenseSurfaceWorkloads,
			Count: totalWorkloads,
		},
	}

	return LicenseBundleInput{
		Name:     LicenseBundleCNSFoundations,
		Surfaces: surfaces,
	}
}

// NewDataIngest30dModuleItem creates a [LicenseModuleItem] for the 30-day Data Ingest module add-on.
func NewDataIngest30dModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngest30d}
}

// NewDataIngest90dModuleItem creates a [LicenseModuleItem] for the 90-day Data Ingest module add-on.
func NewDataIngest90dModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngest90d}
}

// NewDataIngest180dModuleItem creates a [LicenseModuleItem] for the 180-day Data Ingest module add-on.
func NewDataIngest180dModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngest180d}
}

// NewDataIngest365dModuleItem creates a [LicenseModuleItem] for the 365-day Data Ingest module add-on.
func NewDataIngest365dModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngest365d}
}

// NewDataIngestLongRangeRetention1yModuleItem creates a [LicenseModuleItem] for the 1-year Data Ingest Long
// Range Retention module add-on.
func NewDataIngestLongRangeRetention1yModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngestLongRangeRetention1y}
}

// NewDataIngestLongRangeRetention2yModuleItem creates a [LicenseModuleItem] for the 2-year Data Ingest Long
// Range Retention module add-on.
func NewDataIngestLongRangeRetention2yModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngestLongRangeRetention2y}
}

// NewDataIngestLongRangeRetention3yModuleItem creates a [LicenseModuleItem] for the 3-year Data Ingest Long
// Range Retention module add-on.
func NewDataIngestLongRangeRetention3yModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngestLongRangeRetention3y}
}

// NewDataIngestLongRangeRetention4yModuleItem creates a [LicenseModuleItem] for the 4-year Data Ingest Long
// Range Retention module add-on.
func NewDataIngestLongRangeRetention4yModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngestLongRangeRetention4y}
}

// NewDataIngestLongRangeRetention5yModuleItem creates a [LicenseModuleItem] for the 5-year Data Ingest Long
// Range Retention module add-on.
func NewDataIngestLongRangeRetention5yModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngestLongRangeRetention5y}
}

// NewCloudFunnelModuleItem creates a [LicenseModuleItem] for the Cloud Funnel module add-on.
func NewCloudFunnelModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleCloudFunnel}
}

// NewVigilanceMDRModuleItem creates a [LicenseModuleItem] for the Vigilance MDR module add-on.
func NewVigilanceMDRModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleVigilanceMDR}
}

// NewBinaryVaultBenignFilesModuleItem creates a [LicenseModuleItem] for the Binary Vault Benign Files module add-on.
func NewBinaryVaultBenignFilesModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleBinaryVaultBenignFiles}
}

// NewDataIngestLongRangeEndpointAndCWSModuleItem creates a [LicenseModuleItem] for the Data Ingest Long Range
// Endpoint and CWS module add-on.
func NewDataIngestLongRangeEndpointAndCWSModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleDataIngestLongRangeEndpointAndCWS}
}

// NewNetworkDiscoveryModuleItem creates a [LicenseModuleItem] for the Network Discovery module add-on.
func NewNetworkDiscoveryModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleNetworkDiscovery}
}

// NewSingularityMDRModuleItem creates a [LicenseModuleItem] for the Singularity MDR module add-on.
func NewSingularityMDRModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleSingularityMDR}
}

// NewThreatIntelModuleItem creates a [LicenseModuleItem] for the Threat Intel module add-on.
func NewThreatIntelModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleThreatIntel}
}

// NewPurpleAIFoundationsModuleItem creates a [LicenseModuleItem] for the Purple AI Foundations module add-on.
func NewPurpleAIFoundationsModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModulePurpleAIFoundations}
}

// NewPurpleAISocAnalystModuleItem creates a [LicenseModuleItem] for the Purple AI SOC Analyst module add-on.
func NewPurpleAISocAnalystModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModulePurpleAISocAnalyst}
}

// NewWatchTowerModuleItem creates a [LicenseModuleItem] for the WatchTower module add-on.
func NewWatchTowerModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleWatchTower}
}

// NewVulnerabilityManagementModuleItem creates a [LicenseModuleItem] for the Vulnerability Management module add-on.
func NewVulnerabilityManagementModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleVulnerabilityManagement}
}

// NewUnprotectedEndpointDiscoveryModuleItem creates a [LicenseModuleItem] for the Unprotected Endpoint Discovery
// module add-on.
func NewUnprotectedEndpointDiscoveryModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleUnprotectedEndpointDiscovery}
}

// NewRemoteScriptOrchestrationModuleItem creates a [LicenseModuleItem] for the Remote Script Orchestration module
// add-on.
func NewRemoteScriptOrchestrationModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleRemoteScriptOrchestration}
}

// NewWayfinderEliteModuleItem creates a [LicenseModuleItem] for the Wayfinder Elite module add-on.
func NewWayfinderEliteModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleWayfinderElite}
}

// NewWayfinderEssentialsModuleItem creates a [LicenseModuleItem] for the Wayfinder Essentials module add-on.
func NewWayfinderEssentialsModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleWayfinderEssentials}
}

// NewWayfinderThreatHuntingModuleItem creates a [LicenseModuleItem] for the Wayfinder Threat Hunting module add-on.
func NewWayfinderThreatHuntingModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleWayfinderThreatHunting}
}

// NewRemoteOpsForensicsModuleItem creates a [LicenseModuleItem] for the Remote Ops Forensics module add-on.
func NewRemoteOpsForensicsModuleItem() LicenseModuleItem {
	return LicenseModuleItem{Name: LicenseModuleRemoteOpsForensics}
}

// NewXDRDataRetentionSettingInput sets the XDR data retention license setting based on the specified number of days.
//
// Technically only certain discrete values are valid for the API (e.g. 14, 30, 90, 180, 365), but the function
// accepts any positive integer and maps it to the nearest valid maximum value. For example, a value of 45 would map
// to "90 Days", while values of 181 and 400 would map to "365 Days".  Values less than or equal to 14 map to "14
// Days".
func NewXDRDataRetentionSettingInput(days int) LicenseSettingInput {
	var setting string

	switch {
	case days > 180: //nolint:mnd
		setting = "365 Days"
	case days > 90: //nolint:mnd
		setting = "180 Days"
	case days > 30: //nolint:mnd
		setting = "90 Days"
	case days > 14: //nolint:mnd
		setting = "30 Days"
	default:
		setting = "14 Days"
	}

	return LicenseSettingInput{
		GroupName: LicenseSettingXDRDataRetention,
		Setting:   setting,
	}
}

// NewEDRDataRetentionSettingInput sets the EDR data retention license setting.
//
// This value is **always** set to 365 days.
func NewEDRDataRetentionSettingInput() LicenseSettingInput {
	return LicenseSettingInput{
		GroupName: LicenseSettingEDRDataRetention,
		Setting:   "365 Days",
	}
}

// NewRemoteShellSettingInput sets the Remote Shell license setting based on whether the feature should be enabled
// or not.
func NewRemoteShellSettingInput(enabled bool) LicenseSettingInput {
	var setting string
	if enabled {
		setting = "Enabled"
	} else {
		setting = "Disabled"
	}

	return LicenseSettingInput{
		GroupName: LicenseSettingRemoteShell,
		Setting:   setting,
	}
}

// NewMarketplaceAccessSettingInput sets the Marketplace Access license setting based on whether the feature should
// be enabled or not.
func NewMarketplaceAccessSettingInput(enabled bool) LicenseSettingInput {
	var setting string
	if enabled {
		setting = "Available"
	} else {
		setting = "No Access"
	}

	return LicenseSettingInput{
		GroupName: LicenseSettingMarketplaceAccess,
		Setting:   setting,
	}
}

// NewNetworkDiscoveryConsolidationLevelSettingInput sets the Network Discovery Consolidation Level license setting
// based on the specified level, which must be [LicenseSettingNetworkDiscoveryConsolidationLevelAccount] or
// [LicenseSettingNetworkDiscoveryConsolidationLevelSite].
//
// If the given level is invalid, the function defaults to [LicenseSettingNetworkDiscoveryConsolidationLevelAccount].
func NewNetworkDiscoveryConsolidationLevelSettingInput(level string) LicenseSettingInput {
	var setting string

	switch strings.ToLower(level) {
	case "site":
		setting = LicenseSettingNetworkDiscoveryConsolidationLevelSite
	default:
		setting = LicenseSettingNetworkDiscoveryConsolidationLevelAccount
	}

	return LicenseSettingInput{
		GroupName: LicenseSettingNetworkDiscoveryConsolidationLevel,
		Setting:   setting,
	}
}

// NewIdentitySecurityPostureModeSettingInput sets the Identity Security Posture Mode license setting based on the
// specified mode, which must be [LicenseSettingIdentitySecurityPostureModeFull] or
// [LicenseSettingIdentitySecurityPostureModeLite].
//
// If the given mode is invalid, the function defaults to [LicenseSettingIdentitySecurityPostureModeFull].
func NewIdentitySecurityPostureModeSettingInput(mode string) LicenseSettingInput {
	var setting string

	switch strings.ToLower(mode) {
	case "lite":
		setting = LicenseSettingIdentitySecurityPostureModeLite
	default:
		setting = LicenseSettingIdentitySecurityPostureModeFull
	}

	return LicenseSettingInput{
		GroupName: LicenseSettingIdentitySecurityPostureMode,
		Setting:   setting,
	}
}
