// Package types contains all domain/model/response types for the SentinelOne Go SDK.
package types

// AccountSKU represents a deprecated SKU entry on an account.
type AccountSKU struct {
	Type          string `json:"type,omitempty"`
	TotalLicenses int    `json:"totalLicenses,omitempty"`
	Unlimited     bool   `json:"unlimited,omitempty"`
	AgentsInSKU   int    `json:"agentsInSku,omitempty"`
}

// Account represents a SentinelOne management account.
type Account struct {
	ID                  string       `json:"id,omitempty"`
	Name                string       `json:"name,omitempty"`
	IsDefault           bool         `json:"isDefault,omitempty"`
	AccountType         string       `json:"accountType,omitempty"`
	Expiration          *string      `json:"expiration,omitempty"`
	SalesforceID        string       `json:"salesforceId,omitempty"`
	ExternalID          string       `json:"externalId,omitempty"`
	State               string       `json:"state,omitempty"`
	CreatedAt           *string      `json:"createdAt,omitempty"`
	UpdatedAt           *string      `json:"updatedAt,omitempty"`
	UnlimitedExpiration bool         `json:"unlimitedExpiration,omitempty"`
	ActiveAgents        int          `json:"activeAgents,omitempty"`
	TotalLicenses       int          `json:"totalLicenses,omitempty"`
	NumberOfSites       int          `json:"numberOfSites,omitempty"`
	UsageType           string       `json:"usageType,omitempty"`
	BillingMode         string       `json:"billingMode,omitempty"`
	MakeSOCDefaultUI    bool         `json:"makeSocDefaultUi,omitempty"`
	Creator             string       `json:"creator,omitempty"`
	CreatorID           string       `json:"creatorId,omitempty"`
	Licenses            *Licenses    `json:"licenses,omitempty"`
	IRFields            *IRFields    `json:"irFields,omitempty"`
	SKUs                []AccountSKU `json:"skus,omitempty"`

	// Deprecated fields kept for backwards compatibility
	CoreSites           int  `json:"coreSites,omitempty"`
	ControlSites        int  `json:"controlSites,omitempty"`
	CompleteSites       int  `json:"completeSites,omitempty"`
	TotalCore           int  `json:"totalCore,omitempty"`
	TotalControl        int  `json:"totalControl,omitempty"`
	TotalComplete       int  `json:"totalComplete,omitempty"`
	UnlimitedCore       bool `json:"unlimitedCore,omitempty"`
	UnlimitedControl    bool `json:"unlimitedControl,omitempty"`
	UnlimitedComplete   bool `json:"unlimitedComplete,omitempty"`
	AgentsInCoreSKU     int  `json:"agentsInCoreSku,omitempty"`
	AgentsInControlSKU  int  `json:"agentsInControlSku,omitempty"`
	AgentsInCompleteSKU int  `json:"agentsInCompleteSku,omitempty"`
}

// UninstallPasswordMetadata holds metadata about an account's uninstall password.
type UninstallPasswordMetadata struct {
	CreatedAt   *string `json:"createdAt,omitempty"`
	CreatedBy   string  `json:"createdBy,omitempty"`
	CreatedByID string  `json:"createdById,omitempty"`
	ExpiresAt   *string `json:"expiresAt,omitempty"`
	HasPassword bool    `json:"hasPassword,omitempty"`
}

// UninstallPassword holds a generated uninstall password.
type UninstallPassword struct {
	Password string `json:"password,omitempty"`
}
