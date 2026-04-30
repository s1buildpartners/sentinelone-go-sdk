package sentinelone

import (
	"net/url"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// CreateSiteRequest is the request body for POST /sites.
type CreateSiteRequest struct {
	Data CreateSiteData `json:"data"`
}

// CreateSiteData holds the fields for creating a site.
type CreateSiteData struct {
	Name                     string         `json:"name"`
	AccountID                string         `json:"accountId"`
	SiteType                 string         `json:"siteType,omitempty"`
	Expiration               *string        `json:"expiration,omitempty"`
	ExternalID               *string        `json:"externalId,omitempty"`
	Description              string         `json:"description,omitempty"`
	SKU                      string         `json:"sku,omitempty"`
	UnlimitedExpiration      *bool          `json:"unlimitedExpiration,omitempty"`
	UnlimitedLicenses        *bool          `json:"unlimitedLicenses,omitempty"`
	TotalLicenses            *int           `json:"totalLicenses,omitempty"`
	InheritAccountExpiration *bool          `json:"inheritAccountExpiration,omitempty"`
	Inherits                 *bool          `json:"inherits,omitempty"`
	Policy                   *types.Policy  `json:"policy,omitempty"`
	Licenses                 *LicensesInput `json:"licenses,omitempty"`
}

// UpdateSiteRequest is the request body for PUT /sites/{id}.
type UpdateSiteRequest struct {
	Data UpdateSiteData `json:"data"`
}

// UpdateSiteData holds the fields for updating a site.
type UpdateSiteData struct {
	Name                     string         `json:"name,omitempty"`
	SiteType                 string         `json:"siteType,omitempty"`
	Expiration               *string        `json:"expiration,omitempty"`
	ExternalID               *string        `json:"externalId,omitempty"`
	Description              string         `json:"description,omitempty"`
	SKU                      string         `json:"sku,omitempty"`
	UnlimitedExpiration      *bool          `json:"unlimitedExpiration,omitempty"`
	UnlimitedLicenses        *bool          `json:"unlimitedLicenses,omitempty"`
	TotalLicenses            *int           `json:"totalLicenses,omitempty"`
	InheritAccountExpiration *bool          `json:"inheritAccountExpiration,omitempty"`
	Inherits                 *bool          `json:"inherits,omitempty"`
	Policy                   *types.Policy  `json:"policy,omitempty"`
	Licenses                 *LicensesInput `json:"licenses,omitempty"`
}

// BulkUpdateSitesFilter filters which sites to bulk-update.
type BulkUpdateSitesFilter struct {
	AccountIDs []string `json:"accountIds,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	Query      string   `json:"query,omitempty"`
}

// BulkUpdateSitesRequest is the request body for PUT /sites/update-bulk.
type BulkUpdateSitesRequest struct {
	Data   UpdateSiteData        `json:"data"`
	Filter BulkUpdateSitesFilter `json:"filter"`
}

// DuplicateSiteRequest is the request body for POST /sites/duplicate-site.
type DuplicateSiteRequest struct {
	Data DuplicateSiteData `json:"data"`
}

// DuplicateSiteData holds the parameters for duplicating a site.
type DuplicateSiteData struct {
	SiteID     string `json:"siteId"`
	Name       string `json:"name"`
	AccountID  string `json:"accountId,omitempty"`
	CopyPolicy *bool  `json:"copyPolicy,omitempty"`
}

// ReactivateSiteRequest is the request body for PUT /sites/{id}/reactivate.
type ReactivateSiteRequest struct {
	Data ReactivateSiteData `json:"data"`
}

// ReactivateSiteData holds reactivation parameters for a site.
type ReactivateSiteData struct {
	Expiration          *string `json:"expiration,omitempty"`
	UnlimitedExpiration *bool   `json:"unlimitedExpiration,omitempty"`
}

// UpdateLocalAuthorizationRequest is the request body for PUT /sites/{id}/local-authorization.
type UpdateLocalAuthorizationRequest struct {
	SiteAuthorization *string `json:"siteAuthorization,omitempty"`
}

// ListSitesParams contains query parameters for GET /web/api/v2.1/sites.
// All fields are optional; zero values are omitted from the request.
//
//   - SiteIDs: limit results to specific site IDs (max 500).
//   - AccountIDs: limit results to sites in these accounts (max 500).
//   - Query: full-text search on name, account name, and description.
//   - State: "active", "expired", or "deleted".
//   - SiteType: "Trial" or "Paid".
//   - SKU / Module: filter by product SKU or module identifier.
//   - RegistrationToken: find the site that owns this registration token.
//   - AvailableMoveSites: when true, only return sites the caller can move agents to.
//   - NameContains / DescriptionContains / AccountNameContains: partial-match filters.
type ListSitesParams struct {
	ListParams

	SiteIDs             []string
	AccountIDs          []string
	Query               string
	Name                string
	IsDefault           *bool
	HealthStatus        *bool
	SiteType            string // Trial, Paid
	State               string // active, expired, deleted
	States              []string
	StatesNin           []string
	Features            []string
	SKU                 string
	Module              string
	ExternalID          string
	Description         string
	AccountID           string
	Expiration          string
	CreatedAt           string
	UpdatedAt           string
	AdminOnly           *bool
	AvailableMoveSites  *bool
	RegistrationToken   string
	AccountNameContains []string
	NameContains        []string
	DescriptionContains []string
}

func (p *ListSitesParams) values() url.Values {
	vals := p.ListParams.values()
	setStringSlice(vals, "siteIds", p.SiteIDs)
	setStringSlice(vals, "accountIds", p.AccountIDs)
	setString(vals, "query", &p.Query)
	setString(vals, "name", &p.Name)
	setBool(vals, "isDefault", p.IsDefault)
	setBool(vals, "healthStatus", p.HealthStatus)
	setString(vals, "siteType", &p.SiteType)
	setString(vals, "state", &p.State)
	setStringSlice(vals, "states", p.States)
	setStringSlice(vals, "statesNin", p.StatesNin)
	setStringSlice(vals, "features", p.Features)
	setString(vals, "sku", &p.SKU)
	setString(vals, "module", &p.Module)
	setString(vals, "externalId", &p.ExternalID)
	setString(vals, "description", &p.Description)
	setString(vals, "accountId", &p.AccountID)
	setString(vals, "expiration", &p.Expiration)
	setString(vals, "createdAt", &p.CreatedAt)
	setString(vals, "updatedAt", &p.UpdatedAt)
	setBool(vals, "adminOnly", p.AdminOnly)
	setBool(vals, "availableMoveSites", p.AvailableMoveSites)
	setString(vals, "registrationToken", &p.RegistrationToken)
	setStringSlice(vals, "accountName__contains", p.AccountNameContains)
	setStringSlice(vals, "name__contains", p.NameContains)
	setStringSlice(vals, "description__contains", p.DescriptionContains)

	return vals
}
