package types

// LicenseSurface is a surface within a license bundle.
type LicenseSurface struct {
	Count int    `json:"count,omitempty"`
	Name  string `json:"name,omitempty"`
}

// LicenseBundle is a product bundle entry in a license set.
type LicenseBundle struct {
	DisplayName   string           `json:"displayName,omitempty"`
	MajorVersion  int              `json:"majorVersion,omitempty"`
	MinorVersion  int              `json:"minorVersion,omitempty"`
	Name          string           `json:"name,omitempty"`
	Surfaces      []LicenseSurface `json:"surfaces,omitempty"`
	TotalSurfaces int              `json:"totalSurfaces,omitempty"`
}

// LicenseModule is an add-on module in a license set.
type LicenseModule struct {
	DisplayName  string `json:"displayName,omitempty"`
	MajorVersion int    `json:"majorVersion,omitempty"`
	Name         string `json:"name,omitempty"`
}

// LicenseSetting is a deprecated license setting entry.
type LicenseSetting struct {
	DisplayName             string `json:"displayName,omitempty"`
	SettingGroup            string `json:"settingGroup,omitempty"`
	SettingGroupDisplayName string `json:"settingGroupDisplayName,omitempty"`
}

// Licenses holds the complete license information for an account or site.
type Licenses struct {
	Bundles  []LicenseBundle  `json:"bundles,omitempty"`
	Modules  []LicenseModule  `json:"modules,omitempty"`
	Settings []LicenseSetting `json:"settings,omitempty"`
}

// IRFields holds IR (Incident Response) contact and classification data.
type IRFields struct {
	CompanyName                string  `json:"companyName"`
	ContactFirstName           string  `json:"contactFirstName"`
	ContactLastName            string  `json:"contactLastName"`
	ContactEmail               string  `json:"contactEmail"`
	Region                     string  `json:"region,omitempty"`
	Country                    string  `json:"country"`
	City                       *string `json:"city,omitempty"`
	Postal                     *string `json:"postal,omitempty"`
	NumberOfEmployeesEndpoints int     `json:"numberOfEmployeesEndpoints"`
	Industry                   string  `json:"industry"`
}

// PolicyEngines holds the on/off state for each detection engine within a
// [Policy].  Each field accepts one of the string values "protect", "detect",
// or "off" as defined by the SentinelOne API.
type PolicyEngines struct {
	Reputation             string `json:"reputation,omitempty"`
	PreExecution           string `json:"preExecution,omitempty"`
	PreExecutionSuspicious string `json:"preExecutionSuspicious,omitempty"`
	Executables            string `json:"executables,omitempty"`
	DataFiles              string `json:"dataFiles,omitempty"`
	Exploits               string `json:"exploits,omitempty"`
	Penetration            string `json:"penetration,omitempty"`
	PUP                    string `json:"pup,omitempty"`
}

// Policy represents a site or account security policy.
type Policy struct {
	NetworkQuarantineOn  *bool          `json:"networkQuarantineOn,omitempty"`
	AutoImmuneOn         *bool          `json:"autoImmuneOn,omitempty"`
	AutoDecommissionOn   *bool          `json:"autoDecommissionOn,omitempty"`
	IsDefault            *bool          `json:"isDefault,omitempty"`
	ResearchOn           *bool          `json:"researchOn,omitempty"`
	AutoMitigationAction string         `json:"autoMitigationAction,omitempty"`
	AutoDecommissionDays *int           `json:"autoDecommissionDays,omitempty"`
	MitigationMode       string         `json:"mitigationMode,omitempty"`
	CreatedAt            *string        `json:"createdAt,omitempty"`
	AgentNotification    *bool          `json:"agentNotification,omitempty"`
	Engines              *PolicyEngines `json:"engines,omitempty"`
	ScanNewAgents        *bool          `json:"scanNewAgents,omitempty"`
	MonitoringOn         *bool          `json:"monitoringOn,omitempty"`
	InheritedFrom        string         `json:"inheritedFrom,omitempty"`
	UpdatedAt            *string        `json:"updatedAt,omitempty"`
}
