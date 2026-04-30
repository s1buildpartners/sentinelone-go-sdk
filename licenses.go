package sentinelone

import "context"

const (
	licensesRootPath = "/licenses"
)

// LicensesClient provides access to the Licenses API group.
// Access it via [Client.Licenses].
type LicensesClient struct{ c *Client }

// UpdateSitesModules adds or removes add-on license modules for sites
// matching the supplied filter.
//
// API: PUT /web/api/v2.1/licenses/update-sites-modules
// Required permission: Sites.edit
//
// Set req.Data.Operation to "add" or "remove".  Populate req.Data.Modules with
// the module names to affect (e.g. "star", "rso").  Populate req.Filter with
// at least one of SiteIDs or AccountIDs to restrict which sites are modified.
//
// The returned integer is the number of sites that were updated.
func (l *LicensesClient) UpdateSitesModules(
	ctx context.Context,
	req UpdateSitesModulesRequest,
) (int, error) {
	var result struct {
		Affected int `json:"affected"`
	}

	_, err := l.c.put(ctx, licensesRootPath+"/update-sites-modules", req, &result)
	if err != nil {
		return 0, err
	}

	return result.Affected, nil
}
