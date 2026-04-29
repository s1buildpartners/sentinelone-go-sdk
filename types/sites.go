package types

// Site represents a SentinelOne site.
type Site struct {
	ID                       string    `json:"id,omitempty"`
	Name                     string    `json:"name,omitempty"`
	AccountID                string    `json:"accountId,omitempty"`
	AccountName              string    `json:"accountName,omitempty"`
	IsDefault                bool      `json:"isDefault,omitempty"`
	HealthStatus             bool      `json:"healthStatus,omitempty"`
	TotalLicenses            int       `json:"totalLicenses,omitempty"`
	ActiveLicenses           int       `json:"activeLicenses,omitempty"`
	SiteType                 string    `json:"siteType,omitempty"`
	Expiration               *string   `json:"expiration,omitempty"`
	ExternalID               string    `json:"externalId,omitempty"`
	Description              string    `json:"description,omitempty"`
	State                    string    `json:"state,omitempty"`
	Suite                    *string   `json:"suite,omitempty"`
	SKU                      *string   `json:"sku,omitempty"`
	CreatedAt                *string   `json:"createdAt,omitempty"`
	UpdatedAt                *string   `json:"updatedAt,omitempty"`
	Licenses                 *Licenses `json:"licenses,omitempty"`
	UsageType                string    `json:"usageType,omitempty"`
	IRFields                 *IRFields `json:"irFields,omitempty"`
	InheritAccountExpiration bool      `json:"inheritAccountExpiration,omitempty"`
	Creator                  string    `json:"creator,omitempty"`
	CreatorID                string    `json:"creatorId,omitempty"`
	UnlimitedExpiration      bool      `json:"unlimitedExpiration,omitempty"`
	UnlimitedLicenses        bool      `json:"unlimitedLicenses,omitempty"`
}

// SitesResponse is the data envelope returned by the list-sites endpoint.
// It includes aggregate license totals alongside the per-site slice.
type SitesResponse struct {
	AllSites struct {
		ActiveLicenses int `json:"activeLicenses,omitempty"`
		TotalLicenses  int `json:"totalLicenses,omitempty"`
	} `json:"allSites"`
	Sites []Site `json:"sites,omitempty"`
}

// SiteToken holds the current registration token for a site, as returned by
// the read-only token endpoint.  Agents use this token to self-register.
// To rotate the token, use [SitesClient.RegenerateKey] which returns a [SiteKey].
type SiteToken struct {
	Token string `json:"token,omitempty"`
}

// SiteKey holds the new registration token returned after a key rotation via
// [SitesClient.RegenerateKey].  The previous token is immediately invalidated.
type SiteKey struct {
	Token string `json:"token,omitempty"`
}

// LocalAuthorization holds site local upgrade/downgrade authorization info.
type LocalAuthorization struct {
	SiteAuthorization *string `json:"siteAuthorization,omitempty"`
}
