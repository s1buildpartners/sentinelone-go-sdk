package sentinelone

import (
	"context"
	"fmt"
	"net/url"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// SitesClient provides access to the Sites API group.
// Access it via [Client.Sites].
type SitesClient struct{ c *Client }

// List returns a paginated list of sites visible to the authenticated user.
//
// API: GET /web/api/v2.1/sites
// Required permission: Sites.view
//
// The response wraps both an aggregate (AllSites totals) and the per-site
// slice in [types.SitesResponse].  Pass nil for params to use API defaults.
func (s *SitesClient) List(
	ctx context.Context,
	params *ListSitesParams,
) (*types.SitesResponse, *types.Pagination, error) {
	var paramVals url.Values
	if params != nil {
		paramVals = params.values()
	}

	var resp types.SitesResponse

	pag, err := s.c.get(ctx, "/sites", paramVals, &resp)
	if err != nil {
		return nil, nil, err
	}

	return &resp, pag, nil
}

// Get returns the site with the given siteID.
//
// API: GET /web/api/v2.1/sites/{site_id}
// Required permission: Sites.view
func (s *SitesClient) Get(ctx context.Context, siteID string) (*types.Site, error) {
	var site types.Site

	_, err := s.c.get(ctx, "/sites/"+siteID, nil, &site)
	if err != nil {
		return nil, err
	}

	return &site, nil
}

// Create creates a new site inside the specified account.
//
// API: POST /web/api/v2.1/sites
// Required permission: Sites.create
//
// Name and AccountID in req.Data are required.  SiteType ("Trial" or "Paid"),
// SKU, and license counts are optional.  Omit Expiration or set
// UnlimitedExpiration true for sites that should never expire.
func (s *SitesClient) Create(ctx context.Context, req CreateSiteRequest) (*types.Site, error) {
	var site types.Site

	_, err := s.c.post(ctx, "/sites", req, &site)
	if err != nil {
		return nil, err
	}

	return &site, nil
}

// Update updates mutable fields on an existing site.
//
// API: PUT /web/api/v2.1/sites/{site_id}
// Required permission: Sites.update
//
// Only non-zero fields in req.Data are applied.  To update the site's security
// policy in the same call, populate req.Data.Policy.
func (s *SitesClient) Update(ctx context.Context, siteID string, req UpdateSiteRequest) (*types.Site, error) {
	var site types.Site

	_, err := s.c.put(ctx, "/sites/"+siteID, req, &site)
	if err != nil {
		return nil, err
	}

	return &site, nil
}

// Delete permanently deletes a site.  This operation cannot be undone;
// any agents assigned to the site will need to be reassigned.
//
// API: DELETE /web/api/v2.1/sites/{site_id}
// Required permission: Sites.delete
func (s *SitesClient) Delete(ctx context.Context, siteID string) error {
	_, err := s.c.delete(ctx, "/sites/"+siteID, nil, nil)

	return err
}

// GetPolicy returns the active security policy for a site.
// If the site inherits from its account, the inherited values are reflected.
//
// API: GET /web/api/v2.1/sites/{site_id}/policy
// Required permission: Sites.editPolicy
func (s *SitesClient) GetPolicy(ctx context.Context, siteID string) (*types.Policy, error) {
	var policy types.Policy

	_, err := s.c.get(ctx, fmt.Sprintf("/sites/%s/policy", siteID), nil, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// UpdatePolicy replaces the security policy for a site, overriding any
// inherited account-level settings for the fields provided.
//
// API: PUT /web/api/v2.1/sites/{site_id}/policy
// Required permission: Sites.editPolicy
func (s *SitesClient) UpdatePolicy(ctx context.Context, siteID string, req UpdatePolicyRequest) (*types.Policy, error) {
	var policy types.Policy

	_, err := s.c.put(ctx, fmt.Sprintf("/sites/%s/policy", siteID), req, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// RevertPolicy removes any site-level policy overrides, causing the site to
// inherit its account's policy again.
//
// API: PUT /web/api/v2.1/sites/{site_id}/revert-policy
// Required permission: Sites.revertPolicy
func (s *SitesClient) RevertPolicy(ctx context.Context, siteID string) error {
	_, err := s.c.put(ctx, fmt.Sprintf("/sites/%s/revert-policy", siteID), RevertPolicyRequest{}, nil)

	return err
}

// GetToken returns the current registration token for a site.
// Agents use this token to self-register into the site.
//
// API: GET /web/api/v2.1/sites/{site_id}/token
// Required permission: Sites.view
func (s *SitesClient) GetToken(ctx context.Context, siteID string) (*types.SiteToken, error) {
	var token types.SiteToken

	_, err := s.c.get(ctx, fmt.Sprintf("/sites/%s/token", siteID), nil, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// RegenerateKey rotates the site's registration key, invalidating the previous
// token.  Agents already registered are unaffected; new agent registrations
// must use the token from the returned [types.SiteKey].
//
// API: PUT /web/api/v2.1/sites/{site_id}/regenerate-key
// Required permission: Sites.regenerateKey
func (s *SitesClient) RegenerateKey(ctx context.Context, siteID string) (*types.SiteKey, error) {
	var key types.SiteKey

	_, err := s.c.put(ctx, fmt.Sprintf("/sites/%s/regenerate-key", siteID), struct{}{}, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

// Reactivate transitions an expired site back to "active".
// Either req.Data.Expiration or req.Data.UnlimitedExpiration must be set.
//
// API: PUT /web/api/v2.1/sites/{site_id}/reactivate
// Required permission: Sites.reactivate
func (s *SitesClient) Reactivate(ctx context.Context, siteID string, req ReactivateSiteRequest) (*types.Site, error) {
	var site types.Site

	_, err := s.c.put(ctx, fmt.Sprintf("/sites/%s/reactivate", siteID), req, &site)
	if err != nil {
		return nil, err
	}

	return &site, nil
}

// ExpireNow immediately transitions an active site to "expired" without
// waiting for its scheduled expiration date.
//
// API: POST /web/api/v2.1/sites/{site_id}/expire-now
// Required permission: Sites.expire
func (s *SitesClient) ExpireNow(ctx context.Context, siteID string) error {
	_, err := s.c.post(ctx, fmt.Sprintf("/sites/%s/expire-now", siteID), struct{}{}, nil)

	return err
}

// GetLocalAuthorization returns the local upgrade/downgrade authorization
// setting for a site.  When set, agents at this site may upgrade or downgrade
// locally until the authorization expiry date.
//
// API: GET /web/api/v2.1/sites/{site_id}/local-authorization
// Required permission: Sites.localAuthorization
func (s *SitesClient) GetLocalAuthorization(ctx context.Context, siteID string) (*types.LocalAuthorization, error) {
	var auth types.LocalAuthorization

	_, err := s.c.get(ctx, fmt.Sprintf("/sites/%s/local-authorization", siteID), nil, &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

// UpdateLocalAuthorization sets or clears the local upgrade/downgrade
// authorization window for a site.  Set SiteAuthorization to an RFC3339
// timestamp to grant authorization until that date, or to nil to revoke it.
//
// API: PUT /web/api/v2.1/sites/{site_id}/local-authorization
// Required permission: Sites.localAuthorization
func (s *SitesClient) UpdateLocalAuthorization(
	ctx context.Context,
	siteID string,
	req UpdateLocalAuthorizationRequest,
) (*types.LocalAuthorization, error) {
	var auth types.LocalAuthorization

	_, err := s.c.put(ctx, fmt.Sprintf("/sites/%s/local-authorization", siteID), req, &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

// Duplicate creates a new site that is a copy of an existing one.
// The source site ID and a new name are required; AccountID optionally places
// the duplicate in a different account.  Set CopyPolicy true to also copy the
// source site's security policy.
//
// API: POST /web/api/v2.1/sites/duplicate-site
// Required permission: Sites.create
func (s *SitesClient) Duplicate(ctx context.Context, req DuplicateSiteRequest) (*types.Site, error) {
	var site types.Site

	_, err := s.c.post(ctx, "/sites/duplicate-site", req, &site)
	if err != nil {
		return nil, err
	}

	return &site, nil
}

// BulkUpdate applies the same update to all sites that match the supplied
// filter.  Only fields set in req.Data are changed.  Use req.Filter.SiteIDs
// or req.Filter.AccountIDs to narrow the target set.
//
// API: PUT /web/api/v2.1/sites/update-bulk
// Required permission: Sites.update
func (s *SitesClient) BulkUpdate(ctx context.Context, req BulkUpdateSitesRequest) error {
	_, err := s.c.put(ctx, "/sites/update-bulk", req, nil)

	return err
}
