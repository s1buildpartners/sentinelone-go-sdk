package sentinelone

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// ---- AccountsClient ----

func TestAccountsClient_List_Success(t *testing.T) {
	cursor := "next-cursor"
	accounts := []types.Account{{ID: "acc1", Name: "Acme Corp"}}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/accounts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, accounts, &types.Pagination{NextCursor: &cursor, TotalItems: 1})
	})

	result, pag, err := cli.Accounts.List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].ID != "acc1" {
		t.Errorf("unexpected result: %+v", result)
	}
	if pag == nil || pag.TotalItems != 1 || *pag.NextCursor != cursor {
		t.Errorf("unexpected pagination: %+v", pag)
	}
}

func TestAccountsClient_List_WithParams(t *testing.T) {
	var receivedQuery url.Values
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.Query()
		writeJSONEnvelope(w, []types.Account{}, &types.Pagination{TotalItems: 0})
	})

	params := &ListAccountsParams{
		ListParams: ListParams{Limit: IntPtr(5)},
		Name:       "Acme",
		State:      "active",
	}
	_, _, err := cli.Accounts.List(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if receivedQuery.Get("limit") != "5" {
		t.Errorf("expected limit=5, got %q", receivedQuery.Get("limit"))
	}
	if receivedQuery.Get("name") != "Acme" {
		t.Errorf("expected name=Acme, got %q", receivedQuery.Get("name"))
	}
	if receivedQuery.Get("state") != "active" {
		t.Errorf("expected state=active, got %q", receivedQuery.Get("state"))
	}
}

func TestAccountsClient_List_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, []types.APIError{{Code: 401, Message: "unauthorized"}})
	})
	_, _, err := cli.Accounts.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	respErr, ok := AsResponseError(err)
	if !ok || respErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 ResponseError, got %v", err)
	}
}

func TestAccountsClient_Get_Success(t *testing.T) {
	account := types.Account{ID: "acc123", Name: "Test Account", State: "active"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/accounts/acc123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, account, nil)
	})

	result, err := cli.Accounts.Get(context.Background(), "acc123")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "acc123" || result.Name != "Test Account" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestAccountsClient_Get_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusNotFound, []types.APIError{{Message: "not found"}})
	})
	_, err := cli.Accounts.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	respErr, ok := AsResponseError(err)
	if !ok || respErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 ResponseError, got %v", err)
	}
}

func TestAccountsClient_Create_Success(t *testing.T) {
	created := types.Account{ID: "new-acc", Name: "New Account"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/accounts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, created, nil)
	})

	req := CreateAccountRequest{Data: CreateAccountData{Name: "New Account", AccountType: "Trial"}}
	result, err := cli.Accounts.Create(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "new-acc" {
		t.Errorf("unexpected id: %q", result.ID)
	}
}

func TestAccountsClient_Create_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, []types.APIError{{Message: "invalid data"}})
	})
	_, err := cli.Accounts.Create(context.Background(), CreateAccountRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_Update_Success(t *testing.T) {
	updated := types.Account{ID: "acc123", Name: "Updated Name"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/accounts/acc123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, updated, nil)
	})

	req := UpdateAccountRequest{Data: UpdateAccountData{Name: "Updated Name"}}
	result, err := cli.Accounts.Update(context.Background(), "acc123", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Updated Name" {
		t.Errorf("unexpected name: %q", result.Name)
	}
}

func TestAccountsClient_Update_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Accounts.Update(context.Background(), "acc123", UpdateAccountRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_GetPolicy_Success(t *testing.T) {
	policy := types.Policy{MitigationMode: "protect", AutoMitigationAction: "quarantine"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/accounts/acc123/policy" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, policy, nil)
	})

	result, err := cli.Accounts.GetPolicy(context.Background(), "acc123")
	if err != nil {
		t.Fatal(err)
	}
	if result.MitigationMode != "protect" {
		t.Errorf("unexpected policy: %+v", result)
	}
}

func TestAccountsClient_GetPolicy_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Accounts.GetPolicy(context.Background(), "acc123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_UpdatePolicy_Success(t *testing.T) {
	policy := types.Policy{MitigationMode: "detect"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSONEnvelope(w, policy, nil)
	})

	req := UpdatePolicyRequest{Data: types.Policy{MitigationMode: "detect"}}
	result, err := cli.Accounts.UpdatePolicy(context.Background(), "acc123", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.MitigationMode != "detect" {
		t.Errorf("unexpected policy mode: %q", result.MitigationMode)
	}
}

func TestAccountsClient_UpdatePolicy_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Accounts.UpdatePolicy(context.Background(), "acc123", UpdatePolicyRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_RevertPolicy_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/accounts/acc123/revert-policy" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})

	err := cli.Accounts.RevertPolicy(context.Background(), "acc123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountsClient_RevertPolicy_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Accounts.RevertPolicy(context.Background(), "acc123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_Reactivate_Success(t *testing.T) {
	exp := "2025-12-31T00:00:00Z"
	account := types.Account{ID: "acc123", State: "active", Expiration: &exp}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/accounts/acc123/reactivate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, account, nil)
	})

	req := ReactivateAccountRequest{Data: ReactivateAccountData{Expiration: StringPtr("2025-12-31T00:00:00Z")}}
	result, err := cli.Accounts.Reactivate(context.Background(), "acc123", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.State != "active" {
		t.Errorf("unexpected state: %q", result.State)
	}
}

func TestAccountsClient_Reactivate_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Accounts.Reactivate(context.Background(), "acc123", ReactivateAccountRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_ExpireNow_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/accounts/acc123/expire-now" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})

	err := cli.Accounts.ExpireNow(context.Background(), "acc123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountsClient_ExpireNow_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Accounts.ExpireNow(context.Background(), "acc123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_GetUninstallPasswordMetadata_Success(t *testing.T) {
	createdAt := "2024-01-01T00:00:00Z"
	meta := types.UninstallPasswordMetadata{
		CreatedAt:   &createdAt,
		CreatedBy:   "admin",
		HasPassword: true,
	}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/accounts/acc123/uninstall-password/metadata" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, meta, nil)
	})

	result, err := cli.Accounts.GetUninstallPasswordMetadata(context.Background(), "acc123")
	if err != nil {
		t.Fatal(err)
	}
	if !result.HasPassword || result.CreatedBy != "admin" {
		t.Errorf("unexpected metadata: %+v", result)
	}
}

func TestAccountsClient_GetUninstallPasswordMetadata_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Accounts.GetUninstallPasswordMetadata(context.Background(), "acc123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_ViewUninstallPassword_Success(t *testing.T) {
	pass := types.UninstallPassword{Password: "secret123"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/accounts/acc123/uninstall-password/view" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, pass, nil)
	})

	result, err := cli.Accounts.ViewUninstallPassword(context.Background(), "acc123")
	if err != nil {
		t.Fatal(err)
	}
	if result.Password != "secret123" {
		t.Errorf("unexpected password: %q", result.Password)
	}
}

func TestAccountsClient_ViewUninstallPassword_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Accounts.ViewUninstallPassword(context.Background(), "acc123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_GenerateUninstallPassword_Success(t *testing.T) {
	pass := types.UninstallPassword{Password: "newpass456"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/accounts/acc123/uninstall-password/generate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, pass, nil)
	})

	result, err := cli.Accounts.GenerateUninstallPassword(context.Background(), "acc123")
	if err != nil {
		t.Fatal(err)
	}
	if result.Password != "newpass456" {
		t.Errorf("unexpected password: %q", result.Password)
	}
}

func TestAccountsClient_GenerateUninstallPassword_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Accounts.GenerateUninstallPassword(context.Background(), "acc123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsClient_RevokeUninstallPassword_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/accounts/acc123/uninstall-password/revoke" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})

	err := cli.Accounts.RevokeUninstallPassword(context.Background(), "acc123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountsClient_RevokeUninstallPassword_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Accounts.RevokeUninstallPassword(context.Background(), "acc123")
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- ListAccountsParams.values() ----

func TestListAccountsParams_Values_Empty(t *testing.T) {
	p := &ListAccountsParams{}
	vals := p.values()
	if len(vals) != 0 {
		t.Errorf("expected empty values, got %v", vals)
	}
}

func TestListAccountsParams_Values_AllFields(t *testing.T) {
	p := &ListAccountsParams{
		ListParams:   ListParams{Limit: IntPtr(10), Skip: IntPtr(5)},
		IDs:          []string{"id1", "id2"},
		AccountIDs:   []string{"acc1"},
		Query:        "search-term",
		Name:         "Exact Name",
		IsDefault:    BoolPtr(true),
		AccountType:  "Trial",
		State:        "active",
		States:       []string{"active", "expired"},
		StatesNin:    []string{"deleted"},
		Features:     []string{"firewall-control"},
		UsageType:    "customer",
		BillingMode:  "subscription",
		SKU:          "core",
		Module:       "edr",
		Expiration:   "2025-12-31",
		CreatedAt:    "2024-01-01",
		UpdatedAt:    "2024-06-01",
		NameContains: []string{"corp", "inc"},
	}
	vals := p.values()

	checks := map[string]string{
		"limit":          "10",
		"skip":           "5",
		"ids":            "id1,id2",
		"accountIds":     "acc1",
		"query":          "search-term",
		"name":           "Exact Name",
		"isDefault":      "true",
		"accountType":    "Trial",
		"state":          "active",
		"states":         "active,expired",
		"statesNin":      "deleted",
		"features":       "firewall-control",
		"usageType":      "customer",
		"billingMode":    "subscription",
		"sku":            "core",
		"module":         "edr",
		"expiration":     "2025-12-31",
		"createdAt":      "2024-01-01",
		"updatedAt":      "2024-06-01",
		"name__contains": "corp,inc",
	}
	for key, want := range checks {
		if got := vals.Get(key); got != want {
			t.Errorf("param %q: got %q, want %q", key, got, want)
		}
	}
}
