package types

// AgentNetworkInterface represents a network interface on a managed endpoint.
type AgentNetworkInterface struct {
	ID                string   `json:"id,omitempty"`
	Inet              []string `json:"inet,omitempty"`
	Inet6             []string `json:"inet6,omitempty"`
	Name              string   `json:"name,omitempty"`
	Physical          string   `json:"physical,omitempty"`
	GatewayIP         string   `json:"gatewayIp,omitempty"`
	GatewayMacAddress string   `json:"gatewayMacAddress,omitempty"`
}

// AgentLocation holds a location-awareness entry for a managed endpoint.
type AgentLocation struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Scope string `json:"scope,omitempty"`
}

// VssVolume holds VSS shadow-copy volume metrics for a Windows endpoint.
type VssVolume struct {
	DiffAreaName                  string  `json:"diffAreaName,omitempty"`
	DiffAreaCurrentUsedBytes      int     `json:"diffAreaCurrentUsedBytes,omitempty"`
	DiffAreaCurrentAllocatedBytes int     `json:"diffAreaCurrentAllocatedBytes,omitempty"`
	DiffAreaMaxLimitBytes         int     `json:"diffAreaMaxLimitBytes,omitempty"`
	DiffAreaFreePercentage        float64 `json:"diffAreaFreePercentage,omitempty"`
}

// DiskMetrics holds disk usage information for a managed endpoint.
type DiskMetrics struct {
	Path                       string  `json:"path,omitempty"`
	VolumeType                 string  `json:"volumeType,omitempty"`
	TotalNumberOfBytes         int     `json:"totalNumberOfBytes,omitempty"`
	TotalNumberOfFreeBytes     int     `json:"totalNumberOfFreeBytes,omitempty"`
	FreeBytesAvailableToCaller int     `json:"freeBytesAvailableToCaller,omitempty"`
	FreePercentage             float64 `json:"freePercentage,omitempty"`
}

// Agent represents a SentinelOne managed endpoint.
type Agent struct {
	ID                             string                  `json:"id,omitempty"`
	UUID                           string                  `json:"uuid,omitempty"`
	ComputerName                   string                  `json:"computerName,omitempty"`
	AccountID                      string                  `json:"accountId,omitempty"`
	AccountName                    string                  `json:"accountName,omitempty"`
	SiteID                         string                  `json:"siteId,omitempty"`
	SiteName                       string                  `json:"siteName,omitempty"`
	GroupID                        string                  `json:"groupId,omitempty"`
	GroupName                      string                  `json:"groupName,omitempty"`
	GroupIP                        string                  `json:"groupIp,omitempty"`
	AgentVersion                   string                  `json:"agentVersion,omitempty"`
	OSType                         string                  `json:"osType,omitempty"`
	OSName                         string                  `json:"osName,omitempty"`
	OSArch                         string                  `json:"osArch,omitempty"`
	OSRevision                     string                  `json:"osRevision,omitempty"`
	OSStartTime                    *string                 `json:"osStartTime,omitempty"`
	OSUsername                     string                  `json:"osUsername,omitempty"`
	Domain                         string                  `json:"domain,omitempty"`
	MachineType                    string                  `json:"machineType,omitempty"`
	ModelName                      string                  `json:"modelName,omitempty"`
	SerialNumber                   string                  `json:"serialNumber,omitempty"`
	CPUID                          string                  `json:"cpuId,omitempty"`
	CPUCount                       int                     `json:"cpuCount,omitempty"`
	CoreCount                      int                     `json:"coreCount,omitempty"`
	TotalMemory                    int                     `json:"totalMemory,omitempty"`
	NetworkStatus                  string                  `json:"networkStatus,omitempty"`
	LastIPToMgmt                   string                  `json:"lastIpToMgmt,omitempty"`
	ExternalIP                     string                  `json:"externalIp,omitempty"`
	ExternalID                     string                  `json:"externalId,omitempty"`
	MachineSID                     string                  `json:"machineSid,omitempty"`
	InstallerType                  string                  `json:"installerType,omitempty"`
	LicenseKey                     string                  `json:"licenseKey,omitempty"`
	StorageName                    string                  `json:"storageName,omitempty"`
	StorageType                    string                  `json:"storageType,omitempty"`
	LocationType                   string                  `json:"locationType,omitempty"`
	MitigationMode                 string                  `json:"mitigationMode,omitempty"`
	MitigationModeSuspicious       string                  `json:"mitigationModeSuspicious,omitempty"`
	ScanStatus                     string                  `json:"scanStatus,omitempty"`
	DetectionState                 string                  `json:"detectionState,omitempty"`
	OperationalState               string                  `json:"operationalState,omitempty"`
	OperationalStateExpiration     *string                 `json:"operationalStateExpiration,omitempty"`
	RemoteProfilingState           string                  `json:"remoteProfilingState,omitempty"`
	RemoteProfilingStateExpiration *string                 `json:"remoteProfilingStateExpiration,omitempty"`
	ConsoleMigrationStatus         string                  `json:"consoleMigrationStatus,omitempty"`
	RangerStatus                   string                  `json:"rangerStatus,omitempty"`
	RangerVersion                  string                  `json:"rangerVersion,omitempty"`
	AppsVulnerabilityStatus        string                  `json:"appsVulnerabilityStatus,omitempty"`
	LastLoggedInUserName           string                  `json:"lastLoggedInUserName,omitempty"`
	ActiveThreats                  int                     `json:"activeThreats,omitempty"`
	IsActive                       bool                    `json:"isActive,omitempty"`
	IsDecommissioned               bool                    `json:"isDecommissioned,omitempty"`
	IsUninstalled                  bool                    `json:"isUninstalled,omitempty"`
	IsPendingUninstall             bool                    `json:"isPendingUninstall,omitempty"`
	IsUpToDate                     bool                    `json:"isUpToDate,omitempty"`
	IsAdConnector                  bool                    `json:"isAdConnector,omitempty"`
	IsHyperAutomate                bool                    `json:"isHyperAutomate,omitempty"`
	Infected                       bool                    `json:"infected,omitempty"`
	AllowRemoteShell               bool                    `json:"allowRemoteShell,omitempty"`
	InRemoteShellSession           bool                    `json:"inRemoteShellSession,omitempty"`
	HasContainerizedWorkload       bool                    `json:"hasContainerizedWorkload,omitempty"`
	EncryptedApplications          bool                    `json:"encryptedApplications,omitempty"`
	FirewallEnabled                bool                    `json:"firewallEnabled,omitempty"`
	NetworkQuarantineEnabled       bool                    `json:"networkQuarantineEnabled,omitempty"`
	LocationEnabled                bool                    `json:"locationEnabled,omitempty"`
	ThreatRebootRequired           bool                    `json:"threatRebootRequired,omitempty"`
	ShowAlertIcon                  bool                    `json:"showAlertIcon,omitempty"`
	LastActiveDate                 *string                 `json:"lastActiveDate,omitempty"`
	RegisteredAt                   *string                 `json:"registeredAt,omitempty"`
	CreatedAt                      *string                 `json:"createdAt,omitempty"`
	UpdatedAt                      *string                 `json:"updatedAt,omitempty"`
	GroupUpdatedAt                 *string                 `json:"groupUpdatedAt,omitempty"`
	PolicyUpdatedAt                *string                 `json:"policyUpdatedAt,omitempty"`
	ScanStartedAt                  *string                 `json:"scanStartedAt,omitempty"`
	ScanFinishedAt                 *string                 `json:"scanFinishedAt,omitempty"`
	ScanAbortedAt                  *string                 `json:"scanAbortedAt,omitempty"`
	LastSuccessfulScanDate         *string                 `json:"lastSuccessfulScanDate,omitempty"`
	FullDiskScanLastUpdatedAt      *string                 `json:"fullDiskScanLastUpdatedAt,omitempty"`
	FirstFullModeTime              *string                 `json:"firstFullModeTime,omitempty"`
	UserActionsNeeded              []string                `json:"userActionsNeeded,omitempty"`
	MissingPermissions             []string                `json:"missingPermissions,omitempty"`
	ActiveProtection               []string                `json:"activeProtection,omitempty"`
	Locations                      []AgentLocation         `json:"locations,omitempty"`
	NetworkInterfaces              []AgentNetworkInterface `json:"networkInterfaces,omitempty"`
}
