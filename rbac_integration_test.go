//go:build integration

package sentinelone_test

import (
	"context"
	"testing"

	s1 "github.com/s1buildpartners/sentinelone-go-sdk"
)

func TestIntegration_RBAC_List(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()

	roles, pag, err := cli.RBAC.List(ctx, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if pag == nil {
		t.Fatal("expected pagination metadata")
	}

	t.Logf("found %d role(s)", len(roles))
}

func TestIntegration_RBAC_GetTemplate(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()
	accountID := firstAccountID(t, cli)

	template, err := cli.RBAC.GetTemplate(ctx, &s1.GetRoleTemplateParams{
		AccountIDs: []string{accountID},
	})
	if err != nil {
		t.Fatalf("GetTemplate: %v", err)
	}

	if len(template.Pages) == 0 {
		t.Error("expected at least one permission page in template")
	}

	t.Logf("template has %d permission page(s)", len(template.Pages))
}

func TestIntegration_RBAC_CRUD(t *testing.T) {
	cli := integrationClient(t)
	ctx := context.Background()
	accountID := firstAccountID(t, cli)

	// The API requires at least one permission.  Pull the first available
	// identifier from the template so the create call is valid on any console.
	template, err := cli.RBAC.GetTemplate(ctx, &s1.GetRoleTemplateParams{
		AccountIDs: []string{accountID},
	})
	if err != nil {
		t.Fatalf("GetTemplate: %v", err)
	}

	var permIDs []string

outer:
	for _, page := range template.Pages {
		for _, perm := range page.Permissions {
			if perm.Identifier != "" {
				permIDs = append(permIDs, perm.Identifier)
				break outer
			}
		}
	}

	if len(permIDs) == 0 {
		t.Skip("no permission identifiers found in template; skipping")
	}

	name := uniqueName()

	// --- Create ---
	role, err := cli.RBAC.Create(ctx, s1.CreateRoleRequest{
		Data: s1.CreateRoleData{
			Name:          name,
			Description:   "created by go-sdk integration test",
			PermissionIDs: permIDs,
		},
		Filter: s1.RoleScopeFilter{
			AccountIDs: []string{accountID},
		},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if role.ID == "" {
		t.Fatal("Create returned empty ID")
	}

	if role.Name != name {
		t.Errorf("Name = %q, want %q", role.Name, name)
	}

	t.Logf("created role %s (%s)", role.Name, role.ID)

	var deleted bool

	t.Cleanup(func() {
		if deleted {
			return
		}

		if cleanupErr := cli.RBAC.Delete(context.Background(), role.ID); cleanupErr != nil {
			t.Logf("cleanup: delete role %s: %v", role.ID, cleanupErr)
		}
	})

	// --- Get ---
	got, err := cli.RBAC.Get(ctx, role.ID, &s1.GetRolePermissionsParams{
		AccountIDs: []string{accountID},
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.ID != role.ID {
		t.Errorf("Get ID = %q, want %q", got.ID, role.ID)
	}

	// --- Update ---
	const updatedDesc = "updated by go-sdk integration test"

	updated, err := cli.RBAC.Update(ctx, role.ID, s1.UpdateRoleRequest{
		Data: s1.UpdateRoleData{
			Name:          name,
			Description:   updatedDesc,
			PermissionIDs: permIDs,
		},
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if updated.Description != updatedDesc {
		t.Errorf("Description = %q, want %q", updated.Description, updatedDesc)
	}

	// --- Delete ---
	if err = cli.RBAC.Delete(ctx, role.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	deleted = true
	t.Logf("deleted role %s", role.ID)
}
