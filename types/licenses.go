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
