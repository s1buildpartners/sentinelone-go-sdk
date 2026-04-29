package sentinelone

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// ---- SitesClient ----

func TestSitesClient_List_Success(t *testing.T) {
	cursor := "site-cursor"
	sites := types.SitesResponse{
		Sites: []types.Site{{ID: "s1", Name: "Site One"}},
	}
	sites.AllSites.TotalLicenses = 100
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/sites" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, sites, &types.Pagination{NextCursor: &cursor, TotalItems: 1})
	})

	result, pag, err := cli.Sites.List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Sites) != 1 || result.Sites[0].ID != "s1" {
		t.Errorf("unexpected sites: %+v", result.Sites)
	}
	if pag == nil || *pag.NextCursor != cursor {
		t.Errorf("unexpected pagination: %+v", pag)
	}
}

func TestSitesClient_List_WithParams(t *testing.T) {
	var receivedQuery url.Values
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.Query()
		writeJSONEnvelope(w, types.SitesResponse{}, &types.Pagination{})
	})

	params := &ListSitesParams{
		ListParams:  ListParams{Limit: IntPtr(20)},
		State:       "active",
		AccountIDs:  []string{"acc1"},
		SiteType:    "Paid",
		HealthStatus: BoolPtr(true),
	}
	_, _, err := cli.Sites.List(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if receivedQuery.Get("limit") != "20" {
		t.Errorf("expected limit=20")
	}
	if receivedQuery.Get("state") != "active" {
		t.Errorf("expected state=active")
	}
}

func TestSitesClient_List_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, _, err := cli.Sites.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_Get_Success(t *testing.T) {
	site := types.Site{ID: "site123", Name: "My Site", State: "active"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/sites/site123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, site, nil)
	})

	result, err := cli.Sites.Get(context.Background(), "site123")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "site123" || result.State != "active" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestSitesClient_Get_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusNotFound, nil)
	})
	_, err := cli.Sites.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_Create_Success(t *testing.T) {
	created := types.Site{ID: "new-site", Name: "New Site", AccountID: "acc1"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/sites" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, created, nil)
	})

	req := CreateSiteRequest{Data: CreateSiteData{Name: "New Site", AccountID: "acc1"}}
	result, err := cli.Sites.Create(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "new-site" {
		t.Errorf("unexpected id: %q", result.ID)
	}
}

func TestSitesClient_Create_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Sites.Create(context.Background(), CreateSiteRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_Update_Success(t *testing.T) {
	updated := types.Site{ID: "site123", Name: "Updated Site"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSONEnvelope(w, updated, nil)
	})

	req := UpdateSiteRequest{Data: UpdateSiteData{Name: "Updated Site"}}
	result, err := cli.Sites.Update(context.Background(), "site123", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Updated Site" {
		t.Errorf("unexpected name: %q", result.Name)
	}
}

func TestSitesClient_Update_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Sites.Update(context.Background(), "site123", UpdateSiteRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_Delete_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/sites/site123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})

	err := cli.Sites.Delete(context.Background(), "site123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSitesClient_Delete_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Sites.Delete(context.Background(), "site123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_GetPolicy_Success(t *testing.T) {
	policy := types.Policy{MitigationMode: "protect"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/sites/site123/policy" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, policy, nil)
	})

	result, err := cli.Sites.GetPolicy(context.Background(), "site123")
	if err != nil {
		t.Fatal(err)
	}
	if result.MitigationMode != "protect" {
		t.Errorf("unexpected mode: %q", result.MitigationMode)
	}
}

func TestSitesClient_GetPolicy_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Sites.GetPolicy(context.Background(), "site123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_UpdatePolicy_Success(t *testing.T) {
	policy := types.Policy{MitigationMode: "detect"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSONEnvelope(w, policy, nil)
	})

	req := UpdatePolicyRequest{Data: types.Policy{MitigationMode: "detect"}}
	result, err := cli.Sites.UpdatePolicy(context.Background(), "site123", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.MitigationMode != "detect" {
		t.Errorf("unexpected mode: %q", result.MitigationMode)
	}
}

func TestSitesClient_UpdatePolicy_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Sites.UpdatePolicy(context.Background(), "site123", UpdatePolicyRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_RevertPolicy_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/sites/site123/revert-policy" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})

	err := cli.Sites.RevertPolicy(context.Background(), "site123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSitesClient_RevertPolicy_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Sites.RevertPolicy(context.Background(), "site123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_GetToken_Success(t *testing.T) {
	token := types.SiteToken{Token: "registration-token-xyz"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/sites/site123/token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, token, nil)
	})

	result, err := cli.Sites.GetToken(context.Background(), "site123")
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "registration-token-xyz" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestSitesClient_GetToken_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Sites.GetToken(context.Background(), "site123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_RegenerateKey_Success(t *testing.T) {
	key := types.SiteKey{Token: "new-key-abc"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/sites/site123/regenerate-key" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, key, nil)
	})

	result, err := cli.Sites.RegenerateKey(context.Background(), "site123")
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "new-key-abc" {
		t.Errorf("unexpected key: %q", result.Token)
	}
}

func TestSitesClient_RegenerateKey_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Sites.RegenerateKey(context.Background(), "site123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_Reactivate_Success(t *testing.T) {
	site := types.Site{ID: "site123", State: "active"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/sites/site123/reactivate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, site, nil)
	})

	req := ReactivateSiteRequest{Data: ReactivateSiteData{UnlimitedExpiration: BoolPtr(true)}}
	result, err := cli.Sites.Reactivate(context.Background(), "site123", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.State != "active" {
		t.Errorf("unexpected state: %q", result.State)
	}
}

func TestSitesClient_Reactivate_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Sites.Reactivate(context.Background(), "site123", ReactivateSiteRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_ExpireNow_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/sites/site123/expire-now" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})

	err := cli.Sites.ExpireNow(context.Background(), "site123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSitesClient_ExpireNow_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Sites.ExpireNow(context.Background(), "site123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_GetLocalAuthorization_Success(t *testing.T) {
	authStr := "2025-12-31T00:00:00Z"
	auth := types.LocalAuthorization{SiteAuthorization: &authStr}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/sites/site123/local-authorization" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, auth, nil)
	})

	result, err := cli.Sites.GetLocalAuthorization(context.Background(), "site123")
	if err != nil {
		t.Fatal(err)
	}
	if result.SiteAuthorization == nil || *result.SiteAuthorization != authStr {
		t.Errorf("unexpected auth: %+v", result)
	}
}

func TestSitesClient_GetLocalAuthorization_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Sites.GetLocalAuthorization(context.Background(), "site123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_UpdateLocalAuthorization_Success(t *testing.T) {
	authStr := "2026-01-01T00:00:00Z"
	auth := types.LocalAuthorization{SiteAuthorization: &authStr}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSONEnvelope(w, auth, nil)
	})

	req := UpdateLocalAuthorizationRequest{SiteAuthorization: &authStr}
	result, err := cli.Sites.UpdateLocalAuthorization(context.Background(), "site123", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.SiteAuthorization == nil {
		t.Error("expected non-nil authorization")
	}
}

func TestSitesClient_UpdateLocalAuthorization_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Sites.UpdateLocalAuthorization(context.Background(), "site123", UpdateLocalAuthorizationRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_Duplicate_Success(t *testing.T) {
	dup := types.Site{ID: "dup-site", Name: "Duplicate Site"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/sites/duplicate-site" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, dup, nil)
	})

	req := DuplicateSiteRequest{
		Data: DuplicateSiteData{SiteID: "site123", Name: "Duplicate Site", CopyPolicy: BoolPtr(true)},
	}
	result, err := cli.Sites.Duplicate(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "dup-site" {
		t.Errorf("unexpected id: %q", result.ID)
	}
}

func TestSitesClient_Duplicate_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Sites.Duplicate(context.Background(), DuplicateSiteRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSitesClient_BulkUpdate_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/sites/update-bulk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})

	req := BulkUpdateSitesRequest{
		Data:   UpdateSiteData{Description: "Bulk updated"},
		Filter: BulkUpdateSitesFilter{AccountIDs: []string{"acc1"}},
	}
	err := cli.Sites.BulkUpdate(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSitesClient_BulkUpdate_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Sites.BulkUpdate(context.Background(), BulkUpdateSitesRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- ListSitesParams.values() ----

func TestListSitesParams_Values_Empty(t *testing.T) {
	p := &ListSitesParams{}
	vals := p.values()
	if len(vals) != 0 {
		t.Errorf("expected empty values, got %v", vals)
	}
}

func TestListSitesParams_Values_AllFields(t *testing.T) {
	p := &ListSitesParams{
		ListParams:          ListParams{Limit: IntPtr(25)},
		SiteIDs:             []string{"s1", "s2"},
		AccountIDs:          []string{"acc1"},
		Query:               "search",
		Name:                "Site Name",
		IsDefault:           BoolPtr(false),
		HealthStatus:        BoolPtr(true),
		SiteType:            "Paid",
		State:               "active",
		States:              []string{"active", "expired"},
		StatesNin:           []string{"deleted"},
		Features:            []string{"deep-visibility"},
		SKU:                 "complete",
		Module:              "edr",
		ExternalID:          "ext123",
		Description:         "A test site",
		AccountID:           "acc1",
		Expiration:          "2025-01-01",
		CreatedAt:           "2024-01-01",
		UpdatedAt:           "2024-06-01",
		AdminOnly:           BoolPtr(true),
		AvailableMoveSites:  BoolPtr(false),
		RegistrationToken:   "tok123",
		AccountNameContains: []string{"Corp"},
		NameContains:        []string{"site"},
		DescriptionContains: []string{"test"},
	}
	vals := p.values()

	checks := map[string]string{
		"limit":                  "25",
		"siteIds":                "s1,s2",
		"accountIds":             "acc1",
		"query":                  "search",
		"name":                   "Site Name",
		"isDefault":              "false",
		"healthStatus":           "true",
		"siteType":               "Paid",
		"state":                  "active",
		"states":                 "active,expired",
		"statesNin":              "deleted",
		"features":               "deep-visibility",
		"sku":                    "complete",
		"module":                 "edr",
		"externalId":             "ext123",
		"description":            "A test site",
		"accountId":              "acc1",
		"expiration":             "2025-01-01",
		"createdAt":              "2024-01-01",
		"updatedAt":              "2024-06-01",
		"adminOnly":              "true",
		"availableMoveSites":     "false",
		"registrationToken":      "tok123",
		"accountName__contains":  "Corp",
		"name__contains":         "site",
		"description__contains":  "test",
	}
	for key, want := range checks {
		if got := vals.Get(key); got != want {
			t.Errorf("param %q: got %q, want %q", key, got, want)
		}
	}
}

// Suppress unused import warning.
var _ url.Values
