//go:build integration

package sentinelone_test

import (
	"context"
	"testing"

	s1 "github.com/s1buildpartners/sentinelone-go-sdk"
)

func TestIntegration_Sites_List(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	resp, pag, err := cli.Sites.List(ctx, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(resp.Sites) == 0 {
		t.Fatal("expected at least one site")
	}

	if pag == nil {
		t.Fatal("expected pagination metadata")
	}

	t.Logf("found %d site(s); first: %s (%s)", len(resp.Sites), resp.Sites[0].Name, resp.Sites[0].ID)
}

func TestIntegration_Sites_CRUD(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()
	accountID := firstAccountID(t, cli)

	// Borrow the SKU from an existing site — valid values are console-specific
	// and cannot be hard-coded portably.  Always use SiteType "Trial" so the
	// create succeeds on both trial and paid accounts.
	existingResp, _, err := cli.Sites.List(ctx, &s1.ListSitesParams{
		AccountIDs: []string{accountID},
		ListParams: s1.ListParams{Limit: s1.IntPtr(1)},
	})
	if err != nil {
		t.Fatalf("list sites (for SKU): %v", err)
	}

	if len(existingResp.Sites) == 0 {
		t.Skip("no existing sites to borrow SKU from; skipping")
	}

	ref := existingResp.Sites[0]
	if ref.SKU == nil || *ref.SKU == "" {
		t.Skip("first site has no SKU set; skipping")
	}

	sku := *ref.SKU
	name := uniqueName()

	// --- Create ---
	site, err := cli.Sites.Create(ctx, s1.CreateSiteRequest{
		Data: s1.CreateSiteData{
			Name:                name,
			AccountID:           accountID,
			SiteType:            "Trial",
			SKU:                 sku,
			UnlimitedExpiration: s1.BoolPtr(true),
			UnlimitedLicenses:   s1.BoolPtr(true),
		},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if site.ID == "" {
		t.Fatal("Create returned empty ID")
	}

	if site.Name != name {
		t.Errorf("Name = %q, want %q", site.Name, name)
	}

	t.Logf("created site %s (%s)", site.Name, site.ID)

	var deleted bool

	t.Cleanup(func() {
		if deleted {
			return
		}

		if cleanupErr := cli.Sites.Delete(context.Background(), site.ID); cleanupErr != nil {
			t.Logf("cleanup: delete site %s: %v", site.ID, cleanupErr)
		}
	})

	// --- Get ---
	got, err := cli.Sites.Get(ctx, site.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.ID != site.ID {
		t.Errorf("Get ID = %q, want %q", got.ID, site.ID)
	}

	// --- Update ---
	const updatedDesc = "updated by go-sdk integration test"

	updated, err := cli.Sites.Update(ctx, site.ID, s1.UpdateSiteRequest{
		Data: s1.UpdateSiteData{
			Description: updatedDesc,
		},
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if updated.Description != updatedDesc {
		t.Errorf("Description = %q, want %q", updated.Description, updatedDesc)
	}

	// --- GetPolicy ---
	if _, err = cli.Sites.GetPolicy(ctx, site.ID); err != nil {
		t.Errorf("GetPolicy: %v", err)
	}

	// --- GetToken ---
	token, err := cli.Sites.GetToken(ctx, site.ID)
	if err != nil {
		t.Errorf("GetToken: %v", err)
	} else if token.Token == "" {
		t.Error("GetToken returned empty token")
	}

	// --- Delete ---
	if err = cli.Sites.Delete(ctx, site.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	deleted = true
	t.Logf("deleted site %s", site.ID)
}
