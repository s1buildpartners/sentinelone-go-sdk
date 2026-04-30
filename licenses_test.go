package sentinelone

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// ---- LicensesClient ----

func TestLicensesClient_UpdateSitesModules_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/licenses/update-sites-modules" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, map[string]int{"affected": 3}, nil)
	})

	req := UpdateSitesModulesRequest{
		Data: UpdateSitesModulesData{
			Operation: "add",
			Modules:   []LicenseModuleItem{{Name: LicenseModuleWatchTower}},
		},
		Filter: UpdateSitesModulesFilter{
			SiteIDs: []string{"225494730938493804"},
		},
	}
	affected, err := cli.Licenses.UpdateSitesModules(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if affected != 3 {
		t.Errorf("expected 3 affected, got %d", affected)
	}
}

func TestLicensesClient_UpdateSitesModules_RemoveOperation(t *testing.T) {
	var body map[string]interface{}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		writeJSONEnvelope(w, map[string]int{"affected": 1}, nil)
	})

	req := UpdateSitesModulesRequest{
		Data: UpdateSitesModulesData{
			Operation: "remove",
			Modules:   []LicenseModuleItem{{Name: LicenseModuleRemoteScriptOrchestration}},
		},
		Filter: UpdateSitesModulesFilter{
			AccountIDs: []string{"acc123"},
			State:      "active",
		},
	}
	affected, err := cli.Licenses.UpdateSitesModules(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Errorf("expected 1 affected, got %d", affected)
	}

	// Verify the operation was serialized correctly.
	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in request body")
	}
	if data["operation"] != "remove" {
		t.Errorf("expected operation=remove, got %v", data["operation"])
	}
}

func TestLicensesClient_UpdateSitesModules_ZeroAffected(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSONEnvelope(w, map[string]int{"affected": 0}, nil)
	})

	req := UpdateSitesModulesRequest{
		Data:   UpdateSitesModulesData{Operation: "add"},
		Filter: UpdateSitesModulesFilter{SiteIDs: []string{"nonexistent"}},
	}
	affected, err := cli.Licenses.UpdateSitesModules(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if affected != 0 {
		t.Errorf("expected 0 affected, got %d", affected)
	}
}

func TestLicensesClient_UpdateSitesModules_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	req := UpdateSitesModulesRequest{
		Data:   UpdateSitesModulesData{Operation: "add"},
		Filter: UpdateSitesModulesFilter{},
	}
	_, err := cli.Licenses.UpdateSitesModules(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
	respErr, ok := AsResponseError(err)
	if !ok || respErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 ResponseError, got %v", err)
	}
}

func TestLicensesClient_UpdateSitesModules_Unauthorized(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, err := cli.Licenses.UpdateSitesModules(context.Background(), UpdateSitesModulesRequest{
		Data:   UpdateSitesModulesData{Operation: "add"},
		Filter: UpdateSitesModulesFilter{SiteIDs: []string{"s1"}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- AccountsClient.UpdateLicenses ----

func TestAccountsClient_UpdateLicenses_Success(t *testing.T) {
	var body map[string]interface{}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/accounts/acc1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		writeJSONEnvelope(w, types.Account{ID: "acc1"}, nil)
	})

	licenses := LicensesInput{
		Bundles: []LicenseBundleInput{
			{
				Name:     LicenseBundleEndpointSecurityComplete,
				Surfaces: []LicenseSurfaceInput{{Name: LicenseSurfaceTotalAgents, Count: 100}},
			},
		},
		Modules: []LicenseModuleItem{{Name: LicenseModuleNetworkDiscovery}},
		Settings: []LicenseSettingInput{
			{GroupName: LicenseSettingNetworkDiscoveryConsolidationLevel, Setting: LicenseSettingNetworkDiscoveryConsolidationLevelAccount},
		},
	}

	account, err := cli.Accounts.UpdateLicenses(context.Background(), "acc1", licenses)
	if err != nil {
		t.Fatal(err)
	}
	if account.ID != "acc1" {
		t.Errorf("unexpected account id: %s", account.ID)
	}

	// Verify licenses block was sent correctly.
	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data object in request body")
	}
	lic, ok := data["licenses"].(map[string]interface{})
	if !ok {
		t.Fatal("expected licenses in data")
	}
	bundles, ok := lic["bundles"].([]interface{})
	if !ok || len(bundles) != 1 {
		t.Fatalf("expected 1 bundle, got %v", lic["bundles"])
	}
	bundle := bundles[0].(map[string]interface{})
	if bundle["name"] != LicenseBundleEndpointSecurityComplete {
		t.Errorf("expected bundle name %q, got %v", LicenseBundleEndpointSecurityComplete, bundle["name"])
	}
}

func TestAccountsClient_UpdateLicenses_UnlimitedCount(t *testing.T) {
	var body map[string]interface{}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		writeJSONEnvelope(w, types.Account{ID: "acc2"}, nil)
	})

	licenses := LicensesInput{
		Bundles: []LicenseBundleInput{
			{
				Name:     LicenseBundleEndpointSecurityCore,
				Surfaces: []LicenseSurfaceInput{{Name: LicenseSurfaceTotalAgents, Count: LicenseSurfaceUnlimitedCount}},
			},
		},
	}

	_, err := cli.Accounts.UpdateLicenses(context.Background(), "acc2", licenses)
	if err != nil {
		t.Fatal(err)
	}

	// Verify count=-1 is serialized (not omitted).
	data := body["data"].(map[string]interface{})
	lic := data["licenses"].(map[string]interface{})
	bundles := lic["bundles"].([]interface{})
	bundle := bundles[0].(map[string]interface{})
	surfaces := bundle["surfaces"].([]interface{})
	surface := surfaces[0].(map[string]interface{})
	count, ok := surface["count"].(float64)
	if !ok || int(count) != LicenseSurfaceUnlimitedCount {
		t.Errorf("expected count=%d, got %v", LicenseSurfaceUnlimitedCount, surface["count"])
	}
}

func TestAccountsClient_UpdateLicenses_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Accounts.UpdateLicenses(context.Background(), "acc1", LicensesInput{})
	if err == nil {
		t.Fatal("expected error")
	}
	respErr, ok := AsResponseError(err)
	if !ok || respErr.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 ResponseError, got %v", err)
	}
}

// ---- SitesClient.UpdateLicenses ----

func TestSitesClient_UpdateLicenses_Success(t *testing.T) {
	var body map[string]interface{}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/sites/site1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		writeJSONEnvelope(w, types.Site{ID: "site1"}, nil)
	})

	licenses := LicensesInput{
		Bundles: []LicenseBundleInput{
			{
				Name:     LicenseBundleEndpointSecurityControl,
				Surfaces: []LicenseSurfaceInput{{Name: LicenseSurfaceTotalAgents, Count: 50}},
			},
		},
	}

	site, err := cli.Sites.UpdateLicenses(context.Background(), "site1", licenses)
	if err != nil {
		t.Fatal(err)
	}
	if site.ID != "site1" {
		t.Errorf("unexpected site id: %s", site.ID)
	}

	data := body["data"].(map[string]interface{})
	lic, ok := data["licenses"].(map[string]interface{})
	if !ok {
		t.Fatal("expected licenses in data")
	}
	bundles := lic["bundles"].([]interface{})
	if len(bundles) != 1 {
		t.Errorf("expected 1 bundle, got %d", len(bundles))
	}
}

func TestSitesClient_UpdateLicenses_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Sites.UpdateLicenses(context.Background(), "site1", LicensesInput{})
	if err == nil {
		t.Fatal("expected error")
	}
	respErr, ok := AsResponseError(err)
	if !ok || respErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 ResponseError, got %v", err)
	}
}

// ---- request struct serialisation ----

func TestUpdateSitesModulesRequest_Serialisation(t *testing.T) {
	req := UpdateSitesModulesRequest{
		Data: UpdateSitesModulesData{
			Operation: "add",
			Modules: []LicenseModuleItem{
				{Name: LicenseModuleWatchTower},
				{Name: LicenseModuleRemoteScriptOrchestration},
			},
		},
		Filter: UpdateSitesModulesFilter{
			SiteIDs:    []string{"s1", "s2"},
			AccountIDs: []string{"acc1"},
			Query:      "prod",
			Name:       "Production",
			State:      "active",
			SiteType:   "Paid",
		},
	}

	raw, err := jsonRoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}

	data, ok := raw["data"].(map[string]interface{})
	if !ok {
		t.Fatal("missing data object")
	}
	if data["operation"] != "add" {
		t.Errorf("operation = %v", data["operation"])
	}

	filter, ok := raw["filter"].(map[string]interface{})
	if !ok {
		t.Fatal("missing filter object")
	}
	if filter["state"] != "active" {
		t.Errorf("state = %v", filter["state"])
	}
}

func TestLicensesInput_Serialisation(t *testing.T) {
	v := map[string]interface{}{
		"data": UpdateAccountData{
			Licenses: &LicensesInput{
				Bundles: []LicenseBundleInput{
					{
						Name:     LicenseBundleEndpointSecurityComplete,
						Surfaces: []LicenseSurfaceInput{{Name: LicenseSurfaceTotalAgents, Count: 10}},
					},
				},
				Modules: []LicenseModuleItem{{Name: LicenseModuleNetworkDiscovery}},
				Settings: []LicenseSettingInput{
					{GroupName: LicenseSettingNetworkDiscoveryConsolidationLevel, Setting: LicenseSettingNetworkDiscoveryConsolidationLevelAccount},
				},
			},
		},
	}

	raw, err := jsonRoundTrip(v)
	if err != nil {
		t.Fatal(err)
	}

	data := raw["data"].(map[string]interface{})
	lic := data["licenses"].(map[string]interface{})
	bundles := lic["bundles"].([]interface{})
	if len(bundles) != 1 {
		t.Fatalf("expected 1 bundle, got %d", len(bundles))
	}

	bundle := bundles[0].(map[string]interface{})
	if bundle["name"] != LicenseBundleEndpointSecurityComplete {
		t.Errorf("bundle name = %v", bundle["name"])
	}

	surfaces := bundle["surfaces"].([]interface{})
	if len(surfaces) != 1 {
		t.Fatalf("expected 1 surface, got %d", len(surfaces))
	}
	surface := surfaces[0].(map[string]interface{})
	if surface["name"] != LicenseSurfaceTotalAgents {
		t.Errorf("surface name = %v", surface["name"])
	}
	if int(surface["count"].(float64)) != 10 {
		t.Errorf("surface count = %v", surface["count"])
	}

	modules := lic["modules"].([]interface{})
	if len(modules) != 1 {
		t.Fatalf("expected 1 module, got %d", len(modules))
	}
	if modules[0].(map[string]interface{})["name"] != LicenseModuleNetworkDiscovery {
		t.Errorf("module name = %v", modules[0])
	}

	settings := lic["settings"].([]interface{})
	if len(settings) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(settings))
	}
	setting := settings[0].(map[string]interface{})
	if setting["groupName"] != LicenseSettingNetworkDiscoveryConsolidationLevel {
		t.Errorf("groupName = %v", setting["groupName"])
	}
	if setting["setting"] != LicenseSettingNetworkDiscoveryConsolidationLevelAccount {
		t.Errorf("setting = %v", setting["setting"])
	}
}

func TestLicenseSurface_ZeroCount_Serialised(t *testing.T) {
	// Count=0 must be serialized (not omitted), as it is a valid entitlement.
	surface := LicenseSurfaceInput{Name: LicenseSurfaceTotalAgents, Count: 0}
	b, err := json.Marshal(surface)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err = json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	count, exists := m["count"]
	if !exists {
		t.Fatal("count field was omitted; must be present even when 0")
	}
	if int(count.(float64)) != 0 {
		t.Errorf("expected count=0, got %v", count)
	}
}

// jsonRoundTrip marshals v to JSON and back into a generic map.
func jsonRoundTrip(v interface{}) (map[string]interface{}, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var out map[string]interface{}
	if err = json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}

	return out, nil
}
