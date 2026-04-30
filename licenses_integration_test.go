//go:build integration

package sentinelone_test

import (
	"context"
	"net/http"
	"testing"

	s1 "github.com/s1buildpartners/sentinelone-go-sdk"
)

func TestIntegration_Licenses_UpdateSitesModules(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	// Resolve the first accessible site to use as the target.
	resp, _, err := cli.Sites.List(ctx, &s1.ListSitesParams{
		State:      "active",
		ListParams: s1.ListParams{Limit: s1.IntPtr(1)},
	})
	if err != nil {
		t.Fatalf("list sites: %v", err)
	}

	if len(resp.Sites) == 0 {
		t.Skip("no active sites accessible; skipping")
	}

	siteID := resp.Sites[0].ID
	t.Logf("targeting site %s (%s)", resp.Sites[0].Name, siteID)

	// "add" with an empty module list is a safe no-op on most consoles.
	// We validate that the API accepts the request and returns a count.
	affected, err := cli.Licenses.UpdateSitesModules(ctx, s1.UpdateSitesModulesRequest{
		Data: s1.UpdateSitesModulesData{
			Operation: "add",
			Modules:   []s1.LicenseModuleItem{},
		},
		Filter: s1.UpdateSitesModulesFilter{
			SiteIDs: []string{siteID},
		},
	})
	if err != nil {
		// A 403 means the token lacks Sites.edit permission — skip rather than
		// fail so the test suite can still run with read-only credentials.
		if respErr, ok := s1.AsResponseError(err); ok && respErr.StatusCode == http.StatusForbidden {
			t.Skipf("Sites.edit permission not available; skipping: %v", err)
		}

		t.Fatalf("UpdateSitesModules: %v", err)
	}

	t.Logf("UpdateSitesModules(add, empty modules) affected %d site(s)", affected)
}
