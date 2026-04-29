package sentinelone

import (
	"context"
	"fmt"
	"net/url"
)

// -- Site types --

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

// SitesResponse is the data envelope returned by GET /sites (has allSites summary + sites array).
type SitesResponse struct {
	AllSites struct {
		ActiveLicenses int `json:"activeLicenses,omitempty"`
		TotalLicenses  int `json:"totalLicenses,omitempty"`
	} `json:"allSites,omitempty"`
	Sites []Site `json:"sites,omitempty"`
}

// SiteToken holds the registration token for a site.
type SiteToken struct {
	Token string `json:"token,omitempty"`
}

// SiteKey holds the regenerated site key response.
type SiteKey struct {
	Token string `json:"token,omitempty"`
}

// LocalAuthorization holds site local upgrade/downgrade authorization info.
type LocalAuthorization struct {
	SiteAuthorization *string `json:"siteAuthorization,omitempty"`
}

// -- Request types --

// CreateSiteRequest is the request body for POST /sites.
type CreateSiteRequest struct {
	Data CreateSiteData `json:"data"`
}

// CreateSiteData holds the fields for creating a site.
type CreateSiteData struct {
	Name                     string   `json:"name"`
	AccountID                string   `json:"accountId"`
	SiteType                 string   `json:"siteType,omitempty"`
	Expiration               *string  `json:"expiration,omitempty"`
	ExternalID               *string  `json:"externalId,omitempty"`
	Description              string   `json:"description,omitempty"`
	SKU                      string   `json:"sku,omitempty"`
	UnlimitedExpiration      *bool    `json:"unlimitedExpiration,omitempty"`
	UnlimitedLicenses        *bool    `json:"unlimitedLicenses,omitempty"`
	TotalLicenses            *int     `json:"totalLicenses,omitempty"`
	InheritAccountExpiration *bool    `json:"inheritAccountExpiration,omitempty"`
	Inherits                 *bool    `json:"inherits,omitempty"`
	Policy                   *Policy  `json:"policy,omitempty"`
}

// UpdateSiteRequest is the request body for PUT /sites/{id}.
type UpdateSiteRequest struct {
	Data UpdateSiteData `json:"data"`
}

// UpdateSiteData holds the fields for updating a site.
type UpdateSiteData struct {
	Name                     string  `json:"name,omitempty"`
	SiteType                 string  `json:"siteType,omitempty"`
	Expiration               *string `json:"expiration,omitempty"`
	ExternalID               *string `json:"externalId,omitempty"`
	Description              string  `json:"description,omitempty"`
	SKU                      string  `json:"sku,omitempty"`
	UnlimitedExpiration      *bool   `json:"unlimitedExpiration,omitempty"`
	UnlimitedLicenses        *bool   `json:"unlimitedLicenses,omitempty"`
	TotalLicenses            *int    `json:"totalLicenses,omitempty"`
	InheritAccountExpiration *bool   `json:"inheritAccountExpiration,omitempty"`
	Inherits                 *bool   `json:"inherits,omitempty"`
	Policy                   *Policy `json:"policy,omitempty"`
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
	SiteID      string `json:"siteId"`
	Name        string `json:"name"`
	AccountID   string `json:"accountId,omitempty"`
	CopyPolicy  *bool  `json:"copyPolicy,omitempty"`
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

// -- Filter params --

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
	SiteIDs              []string
	AccountIDs           []string
	Query                string
	Name                 string
	IsDefault            *bool
	HealthStatus         *bool
	SiteType             string // Trial, Paid
	State                string // active, expired, deleted
	States               []string
	StatesNin            []string
	Features             []string
	SKU                  string
	Module               string
	ExternalID           string
	Description          string
	AccountID            string
	Expiration           string
	CreatedAt            string
	UpdatedAt            string
	AdminOnly            *bool
	AvailableMoveSites   *bool
	RegistrationToken    string
	AccountNameContains  []string
	NameContains         []string
	DescriptionContains  []string
}

func (p *ListSitesParams) values() url.Values {
	v := p.ListParams.values()
	setStringSlice(v, "siteIds", p.SiteIDs)
	setStringSlice(v, "accountIds", p.AccountIDs)
	setString(v, "query", &p.Query)
	setString(v, "name", &p.Name)
	setBool(v, "isDefault", p.IsDefault)
	setBool(v, "healthStatus", p.HealthStatus)
	setString(v, "siteType", &p.SiteType)
	setString(v, "state", &p.State)
	setStringSlice(v, "states", p.States)
	setStringSlice(v, "statesNin", p.StatesNin)
	setStringSlice(v, "features", p.Features)
	setString(v, "sku", &p.SKU)
	setString(v, "module", &p.Module)
	setString(v, "externalId", &p.ExternalID)
	setString(v, "description", &p.Description)
	setString(v, "accountId", &p.AccountID)
	setString(v, "expiration", &p.Expiration)
	setString(v, "createdAt", &p.CreatedAt)
	setString(v, "updatedAt", &p.UpdatedAt)
	setBool(v, "adminOnly", p.AdminOnly)
	setBool(v, "availableMoveSites", p.AvailableMoveSites)
	setString(v, "registrationToken", &p.RegistrationToken)
	setStringSlice(v, "accountName__contains", p.AccountNameContains)
	setStringSlice(v, "name__contains", p.NameContains)
	setStringSlice(v, "description__contains", p.DescriptionContains)
	return v
}

// -- API methods --

// ListSites returns a paginated list of sites visible to the authenticated user.
//
// API: GET /web/api/v2.1/sites
// Required permission: Sites.view
//
// The response wraps both an aggregate (AllSites totals) and the per-site
// slice in [SitesResponse].  Pass nil for params to use API defaults.
func (c *Client) ListSites(ctx context.Context, params *ListSitesParams) (*SitesResponse, *Pagination, error) {
	var p url.Values
	if params != nil {
		p = params.values()
	}
	var resp SitesResponse
	pag, err := c.get(ctx, "/sites", p, &resp)
	if err != nil {
		return nil, nil, err
	}
	return &resp, pag, nil
}

// GetSite returns the site with the given siteID.
//
// API: GET /web/api/v2.1/sites/{site_id}
// Required permission: Sites.view
func (c *Client) GetSite(ctx context.Context, siteID string) (*Site, error) {
	var site Site
	_, err := c.get(ctx, fmt.Sprintf("/sites/%s", siteID), nil, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// CreateSite creates a new site inside the specified account.
//
// API: POST /web/api/v2.1/sites
// Required permission: Sites.create
//
// Name and AccountID in req.Data are required.  SiteType ("Trial" or "Paid"),
// SKU, and license counts are optional.  Omit Expiration or set
// UnlimitedExpiration true for sites that should never expire.
func (c *Client) CreateSite(ctx context.Context, req CreateSiteRequest) (*Site, error) {
	var site Site
	_, err := c.post(ctx, "/sites", req, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// UpdateSite updates mutable fields on an existing site.
//
// API: PUT /web/api/v2.1/sites/{site_id}
// Required permission: Sites.update
//
// Only non-zero fields in req.Data are applied.  To update the site's security
// policy in the same call, populate req.Data.Policy.
func (c *Client) UpdateSite(ctx context.Context, siteID string, req UpdateSiteRequest) (*Site, error) {
	var site Site
	_, err := c.put(ctx, fmt.Sprintf("/sites/%s", siteID), req, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// DeleteSite permanently deletes a site.  This operation cannot be undone;
// any agents assigned to the site will need to be reassigned.
//
// API: DELETE /web/api/v2.1/sites/{site_id}
// Required permission: Sites.delete
func (c *Client) DeleteSite(ctx context.Context, siteID string) error {
	_, err := c.delete(ctx, fmt.Sprintf("/sites/%s", siteID), nil, nil)
	return err
}

// GetSitePolicy returns the active security policy for a site.
// If the site inherits from its account, the inherited values are reflected.
//
// API: GET /web/api/v2.1/sites/{site_id}/policy
// Required permission: Sites.editPolicy
func (c *Client) GetSitePolicy(ctx context.Context, siteID string) (*Policy, error) {
	var policy Policy
	_, err := c.get(ctx, fmt.Sprintf("/sites/%s/policy", siteID), nil, &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// UpdateSitePolicy replaces the security policy for a site, overriding any
// inherited account-level settings for the fields provided.
//
// API: PUT /web/api/v2.1/sites/{site_id}/policy
// Required permission: Sites.editPolicy
func (c *Client) UpdateSitePolicy(ctx context.Context, siteID string, req UpdatePolicyRequest) (*Policy, error) {
	var policy Policy
	_, err := c.put(ctx, fmt.Sprintf("/sites/%s/policy", siteID), req, &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// RevertSitePolicy removes any site-level policy overrides, causing the site
// to inherit its account's policy again.
//
// API: PUT /web/api/v2.1/sites/{site_id}/revert-policy
// Required permission: Sites.revertPolicy
func (c *Client) RevertSitePolicy(ctx context.Context, siteID string) error {
	_, err := c.put(ctx, fmt.Sprintf("/sites/%s/revert-policy", siteID), RevertPolicyRequest{}, nil)
	return err
}

// GetSiteToken returns the current registration token for a site.
// Agents use this token to self-register into the site.
//
// API: GET /web/api/v2.1/sites/{site_id}/token
// Required permission: Sites.view
func (c *Client) GetSiteToken(ctx context.Context, siteID string) (*SiteToken, error) {
	var token SiteToken
	_, err := c.get(ctx, fmt.Sprintf("/sites/%s/token", siteID), nil, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// RegenerateSiteKey rotates the site's registration key, invalidating the
// previous token.  Agents already registered are unaffected; new agent
// registrations must use the token from the returned [SiteKey].
//
// API: PUT /web/api/v2.1/sites/{site_id}/regenerate-key
// Required permission: Sites.regenerateKey
func (c *Client) RegenerateSiteKey(ctx context.Context, siteID string) (*SiteKey, error) {
	var key SiteKey
	_, err := c.put(ctx, fmt.Sprintf("/sites/%s/regenerate-key", siteID), struct{}{}, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// ReactivateSite transitions an expired site back to "active".
// Either req.Data.Expiration or req.Data.UnlimitedExpiration must be set.
//
// API: PUT /web/api/v2.1/sites/{site_id}/reactivate
// Required permission: Sites.reactivate
func (c *Client) ReactivateSite(ctx context.Context, siteID string, req ReactivateSiteRequest) (*Site, error) {
	var site Site
	_, err := c.put(ctx, fmt.Sprintf("/sites/%s/reactivate", siteID), req, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// ExpireSiteNow immediately transitions an active site to "expired" without
// waiting for its scheduled expiration date.
//
// API: POST /web/api/v2.1/sites/{site_id}/expire-now
// Required permission: Sites.expire
func (c *Client) ExpireSiteNow(ctx context.Context, siteID string) error {
	_, err := c.post(ctx, fmt.Sprintf("/sites/%s/expire-now", siteID), struct{}{}, nil)
	return err
}

// GetSiteLocalAuthorization returns the local upgrade/downgrade authorization
// setting for a site.  When set, agents at this site may upgrade or downgrade
// locally until the authorization expiry date.
//
// API: GET /web/api/v2.1/sites/{site_id}/local-authorization
// Required permission: Sites.localAuthorization
func (c *Client) GetSiteLocalAuthorization(ctx context.Context, siteID string) (*LocalAuthorization, error) {
	var auth LocalAuthorization
	_, err := c.get(ctx, fmt.Sprintf("/sites/%s/local-authorization", siteID), nil, &auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// UpdateSiteLocalAuthorization sets or clears the local upgrade/downgrade
// authorization window for a site.  Set SiteAuthorization to an RFC3339
// timestamp to grant authorization until that date, or to nil to revoke it.
//
// API: PUT /web/api/v2.1/sites/{site_id}/local-authorization
// Required permission: Sites.localAuthorization
func (c *Client) UpdateSiteLocalAuthorization(ctx context.Context, siteID string, req UpdateLocalAuthorizationRequest) (*LocalAuthorization, error) {
	var auth LocalAuthorization
	_, err := c.put(ctx, fmt.Sprintf("/sites/%s/local-authorization", siteID), req, &auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// DuplicateSite creates a new site that is a copy of an existing one.
// The source site ID and a new name are required; AccountID optionally places
// the duplicate in a different account.  Set CopyPolicy true to also copy the
// source site's security policy.
//
// API: POST /web/api/v2.1/sites/duplicate-site
// Required permission: Sites.create
func (c *Client) DuplicateSite(ctx context.Context, req DuplicateSiteRequest) (*Site, error) {
	var site Site
	_, err := c.post(ctx, "/sites/duplicate-site", req, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// BulkUpdateSites applies the same update to all sites that match the
// supplied filter.  Only fields set in req.Data are changed.  Use
// req.Filter.SiteIDs or req.Filter.AccountIDs to narrow the target set.
//
// API: PUT /web/api/v2.1/sites/update-bulk
// Required permission: Sites.update
func (c *Client) BulkUpdateSites(ctx context.Context, req BulkUpdateSitesRequest) error {
	_, err := c.put(ctx, "/sites/update-bulk", req, nil)
	return err
}
