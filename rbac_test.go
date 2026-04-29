package sentinelone

import (
	"context"
	"net/http"
	"testing"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// ---- RBACClient ----

func TestRBACClient_List_Success(t *testing.T) {
	cursor := "rbac-cursor"
	roles := []types.Role{{ID: "role1", Name: "Admin"}}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/rbac/roles" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, roles, &types.Pagination{NextCursor: &cursor, TotalItems: 1})
	})

	result, pag, err := cli.RBAC.List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].ID != "role1" {
		t.Errorf("unexpected result: %+v", result)
	}
	if pag == nil || *pag.NextCursor != cursor {
		t.Errorf("unexpected pagination: %+v", pag)
	}
}

func TestRBACClient_List_WithParams(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") != "MyRole" {
			t.Errorf("expected name=MyRole, got %q", r.URL.Query().Get("name"))
		}
		writeJSONEnvelope(w, []types.Role{}, &types.Pagination{})
	})

	params := &ListRolesParams{Name: "MyRole", PredefinedRole: BoolPtr(false)}
	_, _, err := cli.RBAC.List(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRBACClient_List_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, _, err := cli.RBAC.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRBACClient_GetTemplate_Success(t *testing.T) {
	template := types.RoleWithPermissions{
		Role:  types.Role{ID: "", Name: "template"},
		Pages: []types.RolePage{{Name: "Accounts", Identifier: "accounts"}},
	}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/rbac/role" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, template, nil)
	})

	result, err := cli.RBAC.GetTemplate(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Pages) != 1 || result.Pages[0].Identifier != "accounts" {
		t.Errorf("unexpected template: %+v", result)
	}
}

func TestRBACClient_GetTemplate_WithParams(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tenant") != "true" {
			t.Errorf("expected tenant=true, got %q", r.URL.Query().Get("tenant"))
		}
		writeJSONEnvelope(w, types.RoleWithPermissions{}, nil)
	})

	params := &GetRoleTemplateParams{Tenant: BoolPtr(true), AccountIDs: []string{"acc1"}, SiteIDs: []string{"site1"}, GroupIDs: []string{"g1"}}
	_, err := cli.RBAC.GetTemplate(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRBACClient_GetTemplate_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.RBAC.GetTemplate(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRBACClient_Get_Success(t *testing.T) {
	role := types.RoleWithPermissions{
		Role:  types.Role{ID: "role123", Name: "Viewer"},
		Pages: []types.RolePage{{Name: "Dashboard"}},
	}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/rbac/role/role123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, role, nil)
	})

	result, err := cli.RBAC.Get(context.Background(), "role123", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "role123" || result.Name != "Viewer" {
		t.Errorf("unexpected role: %+v", result)
	}
}

func TestRBACClient_Get_WithParams(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") != "Viewer" {
			t.Errorf("expected name=Viewer")
		}
		writeJSONEnvelope(w, types.RoleWithPermissions{}, nil)
	})

	params := &GetRolePermissionsParams{
		Name:             "Viewer",
		Query:            "q",
		AccountIDs:       []string{"acc1"},
		SiteIDs:          []string{"s1"},
		GroupIDs:         []string{"g1"},
		Tenant:           BoolPtr(true),
		CreatedAtLt:      "2024-01-01",
		CreatedAtGt:      "2023-01-01",
		CreatedAtLte:     "2024-01-02",
		CreatedAtGte:     "2023-01-02",
		CreatedAtBetween: "2023-01-01_2024-01-01",
		UpdatedAtLt:      "2024-06-01",
		UpdatedAtGt:      "2023-06-01",
		UpdatedAtLte:     "2024-06-02",
		UpdatedAtGte:     "2023-06-02",
		UpdatedAtBetween: "2023-06-01_2024-06-01",
	}
	_, err := cli.RBAC.Get(context.Background(), "role123", params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRBACClient_Get_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusNotFound, nil)
	})
	_, err := cli.RBAC.Get(context.Background(), "missing", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRBACClient_Create_Success(t *testing.T) {
	created := types.Role{ID: "new-role", Name: "Custom Role"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/rbac/role" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, created, nil)
	})

	req := CreateRoleRequest{
		Data:   CreateRoleData{Name: "Custom Role", Description: "A custom role", PermissionIDs: []string{"p1"}},
		Filter: RoleScopeFilter{AccountIDs: []string{"acc1"}, Tenant: BoolPtr(false)},
	}
	result, err := cli.RBAC.Create(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "new-role" {
		t.Errorf("unexpected id: %q", result.ID)
	}
}

func TestRBACClient_Create_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.RBAC.Create(context.Background(), CreateRoleRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRBACClient_Update_Success(t *testing.T) {
	updated := types.Role{ID: "role123", Name: "Updated Role"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/rbac/role/role123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, updated, nil)
	})

	req := UpdateRoleRequest{
		Data:   UpdateRoleData{Name: "Updated Role", Description: "Updated", PermissionIDs: []string{"p2"}},
		Filter: &RoleScopeFilter{SiteIDs: []string{"s1"}, GroupIDs: []string{"g1"}},
	}
	result, err := cli.RBAC.Update(context.Background(), "role123", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Updated Role" {
		t.Errorf("unexpected name: %q", result.Name)
	}
}

func TestRBACClient_Update_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.RBAC.Update(context.Background(), "role123", UpdateRoleRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRBACClient_Delete_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/rbac/role/role123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})

	err := cli.RBAC.Delete(context.Background(), "role123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRBACClient_Delete_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.RBAC.Delete(context.Background(), "role123")
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- ListRolesParams.values() ----

func TestListRolesParams_Values_Empty(t *testing.T) {
	p := &ListRolesParams{}
	vals := p.values()
	if len(vals) != 0 {
		t.Errorf("expected empty values, got %v", vals)
	}
}

func TestListRolesParams_Values_AllFields(t *testing.T) {
	p := &ListRolesParams{
		ListParams:       ListParams{Limit: IntPtr(10)},
		AccountIDs:       []string{"acc1"},
		SiteIDs:          []string{"s1"},
		GroupIDs:         []string{"g1"},
		TenancyIDs:       []string{"t1"},
		Tenant:           BoolPtr(true),
		Name:             "Admin",
		Query:            "q",
		IDs:              []string{"id1"},
		Creator:          "alice",
		CreatorID:        "u1",
		UpdatedBy:        "bob",
		UpdatedByID:      "u2",
		Description:      "desc",
		AccountName:      "Acme",
		SiteName:         "HQ",
		IncludeParents:   BoolPtr(true),
		IncludeChildren:  BoolPtr(false),
		PredefinedRole:   BoolPtr(false),
		CreatedAt:        "2024-01-01",
		UpdatedAt:        "2024-06-01",
		CreatedAtLt:      "2024-02-01",
		CreatedAtGt:      "2023-12-01",
		CreatedAtLte:     "2024-02-02",
		CreatedAtGte:     "2023-12-02",
		CreatedAtBetween: "2023-12-01_2024-02-01",
		UpdatedAtLt:      "2024-07-01",
		UpdatedAtGt:      "2024-05-01",
		UpdatedAtLte:     "2024-07-02",
		UpdatedAtGte:     "2024-05-02",
		UpdatedAtBetween: "2024-05-01_2024-07-01",
	}
	vals := p.values()

	checks := map[string]string{
		"limit":                  "10",
		"accountIds":             "acc1",
		"siteIds":                "s1",
		"groupIds":               "g1",
		"tenancyIds":             "t1",
		"tenant":                 "true",
		"name":                   "Admin",
		"query":                  "q",
		"ids":                    "id1",
		"creator":                "alice",
		"creatorId":              "u1",
		"updatedBy":              "bob",
		"updatedById":            "u2",
		"description":            "desc",
		"accountName":            "Acme",
		"siteName":               "HQ",
		"includeParents":         "true",
		"includeChildren":        "false",
		"predefinedRole":         "false",
		"createdAt":              "2024-01-01",
		"updatedAt":              "2024-06-01",
		"createdAt__lt":          "2024-02-01",
		"createdAt__gt":          "2023-12-01",
		"createdAt__lte":         "2024-02-02",
		"createdAt__gte":         "2023-12-02",
		"createdAt__between":     "2023-12-01_2024-02-01",
		"updatedAt__lt":          "2024-07-01",
		"updatedAt__gt":          "2024-05-01",
		"updatedAt__lte":         "2024-07-02",
		"updatedAt__gte":         "2024-05-02",
		"updatedAt__between":     "2024-05-01_2024-07-01",
	}
	for key, want := range checks {
		if got := vals.Get(key); got != want {
			t.Errorf("param %q: got %q, want %q", key, got, want)
		}
	}
}

// ---- GetRoleTemplateParams.values() ----

func TestGetRoleTemplateParams_Values_Empty(t *testing.T) {
	p := &GetRoleTemplateParams{}
	vals := p.values()
	if len(vals) != 0 {
		t.Errorf("expected empty, got %v", vals)
	}
}

func TestGetRoleTemplateParams_Values_AllFields(t *testing.T) {
	p := &GetRoleTemplateParams{
		AccountIDs: []string{"acc1"},
		SiteIDs:    []string{"s1"},
		GroupIDs:   []string{"g1"},
		Tenant:     BoolPtr(true),
	}
	vals := p.values()
	if vals.Get("accountIds") != "acc1" {
		t.Errorf("expected acc1, got %q", vals.Get("accountIds"))
	}
	if vals.Get("tenant") != "true" {
		t.Errorf("expected true, got %q", vals.Get("tenant"))
	}
}

// ---- GetRolePermissionsParams.values() ----

func TestGetRolePermissionsParams_Values_Empty(t *testing.T) {
	p := &GetRolePermissionsParams{}
	vals := p.values()
	if len(vals) != 0 {
		t.Errorf("expected empty, got %v", vals)
	}
}

func TestGetRolePermissionsParams_Values_AllFields(t *testing.T) {
	p := &GetRolePermissionsParams{
		AccountIDs:       []string{"acc1"},
		SiteIDs:          []string{"s1"},
		GroupIDs:         []string{"g1"},
		Tenant:           BoolPtr(false),
		Name:             "Viewer",
		Query:            "search",
		CreatedAtLt:      "2024-01-01",
		CreatedAtGt:      "2023-01-01",
		CreatedAtLte:     "2024-01-02",
		CreatedAtGte:     "2023-01-02",
		CreatedAtBetween: "2023-01-01_2024-01-01",
		UpdatedAtLt:      "2024-06-01",
		UpdatedAtGt:      "2023-06-01",
		UpdatedAtLte:     "2024-06-02",
		UpdatedAtGte:     "2023-06-02",
		UpdatedAtBetween: "2023-06-01_2024-06-01",
	}
	vals := p.values()

	checks := map[string]string{
		"accountIds":         "acc1",
		"siteIds":            "s1",
		"groupIds":           "g1",
		"tenant":             "false",
		"name":               "Viewer",
		"query":              "search",
		"createdAt__lt":      "2024-01-01",
		"createdAt__gt":      "2023-01-01",
		"createdAt__lte":     "2024-01-02",
		"createdAt__gte":     "2023-01-02",
		"createdAt__between": "2023-01-01_2024-01-01",
		"updatedAt__lt":      "2024-06-01",
		"updatedAt__gt":      "2023-06-01",
		"updatedAt__lte":     "2024-06-02",
		"updatedAt__gte":     "2023-06-02",
		"updatedAt__between": "2023-06-01_2024-06-01",
	}
	for key, want := range checks {
		if got := vals.Get(key); got != want {
			t.Errorf("param %q: got %q, want %q", key, got, want)
		}
	}
}
