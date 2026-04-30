package sentinelone

import "net/url"

// ListAgentsParams contains query parameters for GET /web/api/v2.1/agents.
// All fields are optional; zero values are omitted from the request.
//
//   - IDs / IDsNin: include or exclude specific agent IDs.
//   - SiteIDs / AccountIDs / GroupIDs: scope the query to specific console objects.
//   - Query: free-text search across applicable agent attributes.
//   - ComputerName: exact computer name match.
//   - UUID / UUIDs: filter by endpoint UUID.
//   - OSTypes / OSTypesNin: include or exclude OS families ("windows", "linux", "macos").
//   - OSArch / OSArches / OSArchesNin: OS architecture filters ("32 bit", "64 bit", "ARM64").
//   - MachineTypes / MachineTypesNin: machine categories ("desktop", "laptop", "server", …).
//   - NetworkStatuses / NetworkStatusesNin: agent network state filters.
//   - ScanStatuses / ScanStatusesNin: full-disk scan state filters.
//   - AgentVersions / AgentVersionsNin: include or exclude specific agent versions.
//   - MitigationMode / MitigationModeSuspicious: active protection policy modes.
//   - ActiveThreats / ActiveThreatsGt: filter by active threat count.
//   - UserActionsNeeded / UserActionsNeededNin: pending user action filters.
//   - IsActive / IsDecommissioned / IsUninstalled / IsPendingUninstall / IsUpToDate:
//     boolean state filters (pass nil to omit the parameter).
//   - Infected: when true, only return agents with at least one active threat.
//   - RegisteredAtBetween / LastActiveDateBetween / CreatedAtBetween / UpdatedAtBetween:
//     timestamp-range filters in "<from_ms>-<to_ms>" format.
//   - ComputerNameContains / ExternalIPContains: partial-match array filters.
type ListAgentsParams struct {
	ListParams

	IDs        []string
	IDsNin     []string
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string

	Query        string
	ComputerName string
	UUID         string
	UUIDs        []string

	IsActive           *bool
	IsDecommissioned   *bool
	IsUninstalled      *bool
	IsPendingUninstall *bool
	IsUpToDate         *bool
	Infected           *bool

	AgentVersions     []string
	AgentVersionsNin  []string
	RangerVersions    []string
	RangerVersionsNin []string
	RangerStatus      string

	OSArch      string
	OSArches    []string
	OSArchesNin []string
	OSTypes     []string
	OSTypesNin  []string

	MachineTypes    []string
	MachineTypesNin []string

	NetworkStatuses    []string
	NetworkStatusesNin []string
	Domains            []string
	DomainsNin         []string

	ScanStatuses    []string
	ScanStatusesNin []string

	MitigationMode           string
	MitigationModeSuspicious string

	ActiveThreats   *int
	ActiveThreatsGt *int

	ConsoleMigrationStatuses    []string
	ConsoleMigrationStatusesNin []string

	OperationalStates    []string
	OperationalStatesNin []string

	UserActionsNeeded    []string
	UserActionsNeededNin []string

	AppsVulnerabilityStatuses    []string
	AppsVulnerabilityStatusesNin []string

	LocationIDs    []string
	LocationIDsNin []string

	HasLocalConfiguration *bool

	ComputerNameContains []string
	ExternalIPContains   []string

	RegisteredAtBetween   string
	LastActiveDateBetween string
	CreatedAtBetween      string
	UpdatedAtBetween      string

	FilterID string
}

func (p *ListAgentsParams) values() url.Values {
	vals := p.ListParams.values()
	p.setIdentityValues(vals)
	p.setStatusValues(vals)

	return vals
}

func (p *ListAgentsParams) setIdentityValues(vals url.Values) {
	setStringSlice(vals, "ids", p.IDs)
	setStringSlice(vals, "idsNin", p.IDsNin)
	setStringSlice(vals, "siteIds", p.SiteIDs)
	setStringSlice(vals, "accountIds", p.AccountIDs)
	setStringSlice(vals, "groupIds", p.GroupIDs)
	setString(vals, "query", &p.Query)
	setString(vals, "computerName", &p.ComputerName)
	setString(vals, "uuid", &p.UUID)
	setStringSlice(vals, "uuids", p.UUIDs)
	setBool(vals, "isActive", p.IsActive)
	setBool(vals, "isPendingUninstall", p.IsPendingUninstall)
	setBool(vals, "isDecommissioned", p.IsDecommissioned)
	setBool(vals, "isUninstalled", p.IsUninstalled)
	setBool(vals, "isUpToDate", p.IsUpToDate)
	setBool(vals, "infected", p.Infected)
	setStringSlice(vals, "agentVersions", p.AgentVersions)
	setStringSlice(vals, "agentVersionsNin", p.AgentVersionsNin)
	setStringSlice(vals, "rangerVersions", p.RangerVersions)
	setStringSlice(vals, "rangerVersionsNin", p.RangerVersionsNin)
	setString(vals, "rangerStatus", &p.RangerStatus)
	setString(vals, "osArch", &p.OSArch)
	setStringSlice(vals, "osArches", p.OSArches)
	setStringSlice(vals, "osArchesNin", p.OSArchesNin)
	setStringSlice(vals, "osTypes", p.OSTypes)
	setStringSlice(vals, "osTypesNin", p.OSTypesNin)
	setStringSlice(vals, "machineTypes", p.MachineTypes)
	setStringSlice(vals, "machineTypesNin", p.MachineTypesNin)
}

func (p *ListAgentsParams) setStatusValues(vals url.Values) {
	setStringSlice(vals, "networkStatuses", p.NetworkStatuses)
	setStringSlice(vals, "networkStatusesNin", p.NetworkStatusesNin)
	setStringSlice(vals, "domains", p.Domains)
	setStringSlice(vals, "domainsNin", p.DomainsNin)
	setStringSlice(vals, "scanStatuses", p.ScanStatuses)
	setStringSlice(vals, "scanStatusesNin", p.ScanStatusesNin)
	setString(vals, "mitigationMode", &p.MitigationMode)
	setString(vals, "mitigationModeSuspicious", &p.MitigationModeSuspicious)
	setInt(vals, "activeThreats", p.ActiveThreats)
	setInt(vals, "activeThreats__gt", p.ActiveThreatsGt)
	setStringSlice(vals, "consoleMigrationStatuses", p.ConsoleMigrationStatuses)
	setStringSlice(vals, "consoleMigrationStatusesNin", p.ConsoleMigrationStatusesNin)
	setStringSlice(vals, "operationalStates", p.OperationalStates)
	setStringSlice(vals, "operationalStatesNin", p.OperationalStatesNin)
	setStringSlice(vals, "userActionsNeeded", p.UserActionsNeeded)
	setStringSlice(vals, "userActionsNeededNin", p.UserActionsNeededNin)
	setStringSlice(vals, "appsVulnerabilityStatuses", p.AppsVulnerabilityStatuses)
	setStringSlice(vals, "appsVulnerabilityStatusesNin", p.AppsVulnerabilityStatusesNin)
	setStringSlice(vals, "locationIds", p.LocationIDs)
	setStringSlice(vals, "locationIdsNin", p.LocationIDsNin)
	setBool(vals, "hasLocalConfiguration", p.HasLocalConfiguration)
	setStringSlice(vals, "computerName__contains", p.ComputerNameContains)
	setStringSlice(vals, "externalIp__contains", p.ExternalIPContains)
	setString(vals, "registeredAt__between", &p.RegisteredAtBetween)
	setString(vals, "lastActiveDate__between", &p.LastActiveDateBetween)
	setString(vals, "createdAt__between", &p.CreatedAtBetween)
	setString(vals, "updatedAt__between", &p.UpdatedAtBetween)
	setString(vals, "filterId", &p.FilterID)
}
