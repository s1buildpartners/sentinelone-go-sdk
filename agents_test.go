package sentinelone

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// ---- AgentsClient ----

func TestAgentsClient_List_Success(t *testing.T) {
	cursor := "next-cursor"
	agents := []types.Agent{{ID: "ag1", ComputerName: "desktop-01"}}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/agents" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, agents, &types.Pagination{NextCursor: &cursor, TotalItems: 1})
	})

	result, pag, err := cli.Agents.List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].ID != "ag1" {
		t.Errorf("unexpected result: %+v", result)
	}
	if pag == nil || pag.TotalItems != 1 || *pag.NextCursor != cursor {
		t.Errorf("unexpected pagination: %+v", pag)
	}
}

func TestAgentsClient_List_WithParams(t *testing.T) {
	var receivedQuery url.Values
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.Query()
		writeJSONEnvelope(w, []types.Agent{}, &types.Pagination{TotalItems: 0})
	})

	params := &ListAgentsParams{
		ListParams:  ListParams{Limit: IntPtr(5)},
		SiteIDs:     []string{"site1", "site2"},
		OSTypes:     []string{"windows"},
		Infected:    BoolPtr(true),
		MachineTypes: []string{"laptop"},
	}
	_, _, err := cli.Agents.List(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if receivedQuery.Get("limit") != "5" {
		t.Errorf("expected limit=5, got %q", receivedQuery.Get("limit"))
	}
	if receivedQuery.Get("siteIds") != "site1,site2" {
		t.Errorf("expected siteIds, got %q", receivedQuery.Get("siteIds"))
	}
	if receivedQuery.Get("osTypes") != "windows" {
		t.Errorf("expected osTypes=windows, got %q", receivedQuery.Get("osTypes"))
	}
	if receivedQuery.Get("infected") != "true" {
		t.Errorf("expected infected=true, got %q", receivedQuery.Get("infected"))
	}
	if receivedQuery.Get("machineTypes") != "laptop" {
		t.Errorf("expected machineTypes=laptop, got %q", receivedQuery.Get("machineTypes"))
	}
}

func TestAgentsClient_List_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, []types.APIError{{Code: 401, Message: "unauthorized"}})
	})
	_, _, err := cli.Agents.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	respErr, ok := AsResponseError(err)
	if !ok || respErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 ResponseError, got %v", err)
	}
}

func TestAgentsClient_Count_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/web/api/v2.1/agents/count" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, map[string]int{"total": 42}, nil)
	})

	total, err := cli.Agents.Count(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if total != 42 {
		t.Errorf("expected 42, got %d", total)
	}
}

func TestAgentsClient_Count_WithParams(t *testing.T) {
	var receivedQuery url.Values
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.Query()
		writeJSONEnvelope(w, map[string]int{"total": 7}, nil)
	})

	params := &ListAgentsParams{
		OSTypes:  []string{"linux"},
		Infected: BoolPtr(false),
	}
	total, err := cli.Agents.Count(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if total != 7 {
		t.Errorf("expected 7, got %d", total)
	}
	if receivedQuery.Get("osTypes") != "linux" {
		t.Errorf("expected osTypes=linux, got %q", receivedQuery.Get("osTypes"))
	}
}

func TestAgentsClient_Count_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Agents.Count(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	respErr, ok := AsResponseError(err)
	if !ok || respErr.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 ResponseError, got %v", err)
	}
}

// ---- ListAgentsParams.values() ----

func TestListAgentsParams_Values_Empty(t *testing.T) {
	params := &ListAgentsParams{}
	vals := params.values()
	if len(vals) != 0 {
		t.Errorf("expected empty values, got %v", vals)
	}
}

func TestListAgentsParams_Values_AllFields(t *testing.T) {
	params := &ListAgentsParams{
		ListParams:                   ListParams{Limit: IntPtr(10), Skip: IntPtr(0)},
		IDs:                          []string{"ag1", "ag2"},
		IDsNin:                       []string{"ag3"},
		SiteIDs:                      []string{"s1"},
		AccountIDs:                   []string{"acc1"},
		GroupIDs:                     []string{"g1"},
		Query:                        "desktop",
		ComputerName:                 "workstation-01",
		UUID:                         "abc-123",
		UUIDs:                        []string{"abc-123", "def-456"},
		IsActive:                     BoolPtr(true),
		IsPendingUninstall:           BoolPtr(false),
		IsDecommissioned:             BoolPtr(false),
		IsUninstalled:                BoolPtr(false),
		IsUpToDate:                   BoolPtr(true),
		Infected:                     BoolPtr(false),
		AgentVersions:                []string{"22.1.0"},
		AgentVersionsNin:             []string{"21.0.0"},
		RangerVersions:               []string{"21.11.0"},
		RangerVersionsNin:            []string{"20.0.0"},
		RangerStatus:                 "Enabled",
		OSArch:                       "64 bit",
		OSArches:                     []string{"64 bit"},
		OSArchesNin:                  []string{"32 bit"},
		OSTypes:                      []string{"windows", "linux"},
		OSTypesNin:                   []string{"macos"},
		MachineTypes:                 []string{"laptop"},
		MachineTypesNin:              []string{"server"},
		NetworkStatuses:              []string{"connected"},
		NetworkStatusesNin:           []string{"disconnected"},
		Domains:                      []string{"corp.local"},
		DomainsNin:                   []string{"workgroup"},
		ScanStatuses:                 []string{"finished"},
		ScanStatusesNin:              []string{"aborted"},
		MitigationMode:               "protect",
		MitigationModeSuspicious:     "detect",
		ActiveThreats:                IntPtr(0),
		ActiveThreatsGt:              IntPtr(1),
		ConsoleMigrationStatuses:     []string{"N/A"},
		ConsoleMigrationStatusesNin:  []string{"Failed"},
		OperationalStates:            []string{"na"},
		OperationalStatesNin:         []string{"ghost_mode"},
		UserActionsNeeded:            []string{"reboot_needed"},
		UserActionsNeededNin:         []string{"none"},
		AppsVulnerabilityStatuses:    []string{"patch_required"},
		AppsVulnerabilityStatusesNin: []string{"up_to_date"},
		LocationIDs:                  []string{"loc1"},
		LocationIDsNin:               []string{"loc2"},
		HasLocalConfiguration:        BoolPtr(false),
		ComputerNameContains:         []string{"desk"},
		ExternalIPContains:           []string{"10.0"},
		RegisteredAtBetween:          "1514978764288-1514978999999",
		LastActiveDateBetween:        "1514978764288-1514978999999",
		CreatedAtBetween:             "1514978764288-1514978999999",
		UpdatedAtBetween:             "1514978764288-1514978999999",
		FilterID:                     "filter-abc",
	}

	vals := params.values()

	checks := map[string]string{
		"limit":                        "10",
		"skip":                         "0",
		"ids":                          "ag1,ag2",
		"idsNin":                       "ag3",
		"siteIds":                      "s1",
		"accountIds":                   "acc1",
		"groupIds":                     "g1",
		"query":                        "desktop",
		"computerName":                 "workstation-01",
		"uuid":                         "abc-123",
		"uuids":                        "abc-123,def-456",
		"isActive":                     "true",
		"isPendingUninstall":           "false",
		"isDecommissioned":             "false",
		"isUninstalled":                "false",
		"isUpToDate":                   "true",
		"infected":                     "false",
		"agentVersions":                "22.1.0",
		"agentVersionsNin":             "21.0.0",
		"rangerVersions":               "21.11.0",
		"rangerVersionsNin":            "20.0.0",
		"rangerStatus":                 "Enabled",
		"osArch":                       "64 bit",
		"osArches":                     "64 bit",
		"osArchesNin":                  "32 bit",
		"osTypes":                      "windows,linux",
		"osTypesNin":                   "macos",
		"machineTypes":                 "laptop",
		"machineTypesNin":              "server",
		"networkStatuses":              "connected",
		"networkStatusesNin":           "disconnected",
		"domains":                      "corp.local",
		"domainsNin":                   "workgroup",
		"scanStatuses":                 "finished",
		"scanStatusesNin":              "aborted",
		"mitigationMode":               "protect",
		"mitigationModeSuspicious":     "detect",
		"activeThreats":                "0",
		"activeThreats__gt":            "1",
		"consoleMigrationStatuses":     "N/A",
		"consoleMigrationStatusesNin":  "Failed",
		"operationalStates":            "na",
		"operationalStatesNin":         "ghost_mode",
		"userActionsNeeded":            "reboot_needed",
		"userActionsNeededNin":         "none",
		"appsVulnerabilityStatuses":    "patch_required",
		"appsVulnerabilityStatusesNin": "up_to_date",
		"locationIds":                  "loc1",
		"locationIdsNin":               "loc2",
		"hasLocalConfiguration":        "false",
		"computerName__contains":       "desk",
		"externalIp__contains":         "10.0",
		"registeredAt__between":        "1514978764288-1514978999999",
		"lastActiveDate__between":      "1514978764288-1514978999999",
		"createdAt__between":           "1514978764288-1514978999999",
		"updatedAt__between":           "1514978764288-1514978999999",
		"filterId":                     "filter-abc",
	}

	for key, want := range checks {
		if got := vals.Get(key); got != want {
			t.Errorf("param %q: got %q, want %q", key, got, want)
		}
	}
}
