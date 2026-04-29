package sentinelone

import (
	"context"
	"net/http"
	"testing"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// ---- UsersClient CRUD ----

func TestUsersClient_List_Success(t *testing.T) {
	cursor := "user-cursor"
	email := "alice@example.com"
	name := "Alice"
	users := []types.User{{ID: "u1", Email: &email, FullName: &name}}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, users, &types.Pagination{NextCursor: &cursor, TotalItems: 1})
	})

	result, pag, err := cli.Users.List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].ID != "u1" {
		t.Errorf("unexpected users: %+v", result)
	}
	if pag == nil || *pag.NextCursor != cursor {
		t.Errorf("unexpected pagination: %+v", pag)
	}
}

func TestUsersClient_List_WithParams(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("email") != "alice@example.com" {
			t.Errorf("expected email param")
		}
		writeJSONEnvelope(w, []types.User{}, &types.Pagination{})
	})
	params := &ListUsersParams{Email: "alice@example.com"}
	_, _, err := cli.Users.List(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_List_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, _, err := cli.Users.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_Get_Success(t *testing.T) {
	email := "bob@example.com"
	user := types.User{ID: "u2", Email: &email, Scope: "account"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/u2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, user, nil)
	})
	result, err := cli.Users.Get(context.Background(), "u2")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "u2" {
		t.Errorf("unexpected id: %q", result.ID)
	}
}

func TestUsersClient_Get_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusNotFound, nil)
	})
	_, err := cli.Users.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_Create_Success(t *testing.T) {
	email := "carol@example.com"
	created := types.User{ID: "u3", Email: &email, Scope: "site"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST")
		}
		if r.URL.Path != "/web/api/v2.1/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, created, nil)
	})
	req := CreateUserRequest{Data: CreateUserData{Email: "carol@example.com", FullName: "Carol", Scope: "site"}}
	result, err := cli.Users.Create(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "u3" {
		t.Errorf("unexpected id: %q", result.ID)
	}
}

func TestUsersClient_Create_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Users.Create(context.Background(), CreateUserRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_Update_Success(t *testing.T) {
	name := "Updated Name"
	updated := types.User{ID: "u1", FullName: &name, Scope: "account"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT")
		}
		writeJSONEnvelope(w, updated, nil)
	})
	req := UpdateUserRequest{Data: UpdateUserData{Scope: "account", FullName: "Updated Name"}}
	result, err := cli.Users.Update(context.Background(), "u1", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.FullName == nil || *result.FullName != "Updated Name" {
		t.Errorf("unexpected name: %v", result.FullName)
	}
}

func TestUsersClient_Update_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Users.Update(context.Background(), "u1", UpdateUserRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_Delete_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE")
		}
		if r.URL.Path != "/web/api/v2.1/users/u1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	err := cli.Users.Delete(context.Background(), "u1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_Delete_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.Delete(context.Background(), "u1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_BulkDelete_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/delete-users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	req := BulkUsersActionRequest{Filter: BulkUsersFilter{IDs: []string{"u1", "u2"}}}
	err := cli.Users.BulkDelete(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_BulkDelete_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.BulkDelete(context.Background(), BulkUsersActionRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- API token management ----

func TestUsersClient_GetAPITokenDetails_Success(t *testing.T) {
	detail := types.APITokenDetail{CreatedAt: "2024-01-01", ExpiresAt: "2025-01-01"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/u1/api-token-details" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, detail, nil)
	})
	result, err := cli.Users.GetAPITokenDetails(context.Background(), "u1")
	if err != nil {
		t.Fatal(err)
	}
	if result.CreatedAt != "2024-01-01" {
		t.Errorf("unexpected createdAt: %q", result.CreatedAt)
	}
}

func TestUsersClient_GetAPITokenDetails_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Users.GetAPITokenDetails(context.Background(), "u1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_GetAPITokenDetailsByToken_Success(t *testing.T) {
	detail := types.APITokenDetail{ExpiresAt: "2025-01-01"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/api-token-details" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, detail, nil)
	})
	req := GetAPITokenDetailsRequest{Data: GetAPITokenDetailsData{Token: "mytoken"}}
	result, err := cli.Users.GetAPITokenDetailsByToken(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.ExpiresAt != "2025-01-01" {
		t.Errorf("unexpected expiresAt: %q", result.ExpiresAt)
	}
}

func TestUsersClient_GetAPITokenDetailsByToken_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Users.GetAPITokenDetailsByToken(context.Background(), GetAPITokenDetailsRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_GenerateAPIToken_Success(t *testing.T) {
	token := types.APITokenResponse{Token: "generated-token"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/generate-api-token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, token, nil)
	})
	req := GenerateAPITokenRequest{Data: GenerateAPITokenData{ForceLegacy: BoolPtr(false)}}
	result, err := cli.Users.GenerateAPIToken(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "generated-token" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_GenerateAPIToken_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Users.GenerateAPIToken(context.Background(), GenerateAPITokenRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_RevokeAPIToken_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/revoke-api-token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	err := cli.Users.RevokeAPIToken(context.Background(), UserIDRequest{Data: UserIDData{UserID: "u1"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_RevokeAPIToken_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.RevokeAPIToken(context.Background(), UserIDRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- 2FA management ----

func TestUsersClient_Enable2FA_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/2fa/enable" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	err := cli.Users.Enable2FA(context.Background(), UserIDRequest{Data: UserIDData{UserID: "u1"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_Enable2FA_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.Enable2FA(context.Background(), UserIDRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_Disable2FA_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/2fa/disable" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	err := cli.Users.Disable2FA(context.Background(), UserIDRequest{Data: UserIDData{UserID: "u1"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_Disable2FA_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.Disable2FA(context.Background(), UserIDRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_Enroll2FA_Success(t *testing.T) {
	resp := types.EnrollTFAResponse{Secret: "TOTP-SECRET", QRCode: "data:image/png;base64,..."}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/enroll-2fa" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, resp, nil)
	})
	result, err := cli.Users.Enroll2FA(context.Background(), UserIDsRequest{Data: UserIDsData{UserIDs: []string{"u1"}}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Secret != "TOTP-SECRET" {
		t.Errorf("unexpected secret: %q", result.Secret)
	}
}

func TestUsersClient_Enroll2FA_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Users.Enroll2FA(context.Background(), UserIDsRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_Reset2FA_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/reset-2fa" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	err := cli.Users.Reset2FA(context.Background(), ResetTFARequest{Data: ResetTFAData{UserID: "u1"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_Reset2FA_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.Reset2FA(context.Background(), ResetTFARequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_Delete2FA_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/delete-2fa" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	err := cli.Users.Delete2FA(context.Background(), DeleteTFARequest{Data: DeleteTFAData{UserID: "u1"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_Delete2FA_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.Delete2FA(context.Background(), DeleteTFARequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- Password management ----

func TestUsersClient_ChangePassword_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/change-password" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	req := ChangePasswordRequest{Data: ChangePasswordData{CurrentPassword: "old", NewPassword: "new"}}
	err := cli.Users.ChangePassword(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_ChangePassword_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	err := cli.Users.ChangePassword(context.Background(), ChangePasswordRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_SendResetPasswordEmail_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/login/send-reset-password-email" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	req := SendResetPasswordRequest{Filter: BulkUsersFilter{Email: "alice@example.com"}}
	err := cli.Users.SendResetPasswordEmail(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_SendResetPasswordEmail_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.SendResetPasswordEmail(context.Background(), SendResetPasswordRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_ForceResetPasswordOnLogin_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/login/force-reset-password-on-login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	req := ForceResetPasswordRequest{Filter: BulkUsersFilter{IDs: []string{"u1"}}}
	err := cli.Users.ForceResetPasswordOnLogin(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_ForceResetPasswordOnLogin_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.ForceResetPasswordOnLogin(context.Background(), ForceResetPasswordRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- App / iframe integration ----

func TestUsersClient_EnableApp_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/enable-app" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	req := EnableAppRequest{Data: EnableAppData{AppID: "app1", Enable: true}}
	err := cli.Users.EnableApp(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_EnableApp_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	err := cli.Users.EnableApp(context.Background(), EnableAppRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_RequestApp_Success(t *testing.T) {
	resp := types.RequestAppResponse{URL: "https://app.example.com/access"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/request-app" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, resp, nil)
	})
	result, err := cli.Users.RequestApp(context.Background(), RequestAppRequest{CurrentPassword: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	if result.URL != "https://app.example.com/access" {
		t.Errorf("unexpected URL: %q", result.URL)
	}
}

func TestUsersClient_RequestApp_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Users.RequestApp(context.Background(), RequestAppRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_GenerateIFrameToken_Success(t *testing.T) {
	token := types.IFrameTokenResponse{Token: "iframe-token-xyz"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/generate-iframe-token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, token, nil)
	})
	req := IFrameUserRequest{Data: IFrameUserData{AccountID: "acc1", SiteID: "s1", Expiration: "2025-01-01"}}
	result, err := cli.Users.GenerateIFrameToken(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "iframe-token-xyz" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_GenerateIFrameToken_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Users.GenerateIFrameToken(context.Background(), IFrameUserRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- Authentication ----

func TestUsersClient_Login_Success(t *testing.T) {
	resp := types.LoginResponse{Token: "session-token", Status: "active"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, resp, nil)
	})
	req := LoginRequest{Username: "alice@example.com", Password: "pass", RememberMe: BoolPtr(true)}
	result, err := cli.Users.Login(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "session-token" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_Login_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, err := cli.Users.Login(context.Background(), LoginRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_LoginContinue_Success(t *testing.T) {
	resp := types.LoginContinueResponse{Token: "full-token", Status: "active"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/login-continue" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, resp, nil)
	})
	req := LoginContinueRequest{Data: LoginContinueData{Token: "temp-token", Code: "123456", Method: "totp"}}
	result, err := cli.Users.LoginContinue(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "full-token" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_LoginContinue_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, err := cli.Users.LoginContinue(context.Background(), LoginContinueRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_LoginByToken_Success(t *testing.T) {
	resp := types.LoginResponse{Token: "session-from-token"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/login/by-token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("token") != "invite-token" {
			t.Errorf("expected token param")
		}
		writeJSONEnvelope(w, resp, nil)
	})
	result, err := cli.Users.LoginByToken(context.Background(), "invite-token")
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "session-from-token" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_LoginByToken_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, err := cli.Users.LoginByToken(context.Background(), "bad-token")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_LoginByAPIToken_Success(t *testing.T) {
	resp := types.LoginResponse{Token: "session-token"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/login/by-api-token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, resp, nil)
	})
	req := LoginByAPITokenRequest{Data: LoginByAPITokenData{APIToken: "api-token"}}
	result, err := cli.Users.LoginByAPIToken(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "session-token" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_LoginByAPIToken_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, err := cli.Users.LoginByAPIToken(context.Background(), LoginByAPITokenRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_LoginSSO_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/login/sso-saml2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, map[string]string{"redirectUrl": "https://idp.example.com/sso"}, nil)
	})
	url, err := cli.Users.LoginSSO(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://idp.example.com/sso" {
		t.Errorf("unexpected URL: %q", url)
	}
}

func TestUsersClient_LoginSSO_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusServiceUnavailable, nil)
	})
	_, err := cli.Users.LoginSSO(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_LoginSSOForScope_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/login/sso-saml2/scope123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, map[string]string{"redirectUrl": "https://idp.example.com/sso/scope"}, nil)
	})
	url, err := cli.Users.LoginSSOForScope(context.Background(), "scope123")
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://idp.example.com/sso/scope" {
		t.Errorf("unexpected URL: %q", url)
	}
}

func TestUsersClient_LoginSSOForScope_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusServiceUnavailable, nil)
	})
	_, err := cli.Users.LoginSSOForScope(context.Background(), "scope123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_SetPassword_Success(t *testing.T) {
	resp := types.SetPasswordResponse{Token: "new-session"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/login/set-password" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, resp, nil)
	})
	req := SetPasswordRequest{Data: SetPasswordData{Token: "reset-token", Password: "newpass"}}
	result, err := cli.Users.SetPassword(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "new-session" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_SetPassword_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Users.SetPassword(context.Background(), SetPasswordRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_Logout_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/logout" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	err := cli.Users.Logout(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_Logout_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	err := cli.Users.Logout(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- Session management ----

func TestUsersClient_ElevateSession_Success(t *testing.T) {
	resp := types.ElevateSessionResponse{Token: "elevated-token"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/auth/elevate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, resp, nil)
	})
	req := ElevateSessionRequest{Data: ElevateSessionData{Password: "mypass"}}
	result, err := cli.Users.ElevateSession(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "elevated-token" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_ElevateSession_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, err := cli.Users.ElevateSession(context.Background(), ElevateSessionRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_SSOReAuth_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/sso-saml2/re-auth" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, map[string]string{"redirectUrl": "https://idp.example.com/reauth"}, nil)
	})
	url, err := cli.Users.SSOReAuth(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://idp.example.com/reauth" {
		t.Errorf("unexpected URL: %q", url)
	}
}

func TestUsersClient_SSOReAuth_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, err := cli.Users.SSOReAuth(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_AuthEULA_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/auth/eula" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	err := cli.Users.AuthEULA(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_AuthEULA_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	err := cli.Users.AuthEULA(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_AuthApp_Success(t *testing.T) {
	resp := types.LoginResponse{Token: "app-session"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/auth/app" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, resp, nil)
	})
	req := AuthAppRequest{Data: AuthAppData{Code: "auth-code"}}
	result, err := cli.Users.AuthApp(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "app-session" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_AuthApp_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusUnauthorized, nil)
	})
	_, err := cli.Users.AuthApp(context.Background(), AuthAppRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_IsTenantAdmin_True(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/tenant-admin-auth-check" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, map[string]bool{"isAdmin": true}, nil)
	})
	isAdmin, err := cli.Users.IsTenantAdmin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !isAdmin {
		t.Error("expected isAdmin=true")
	}
}

func TestUsersClient_IsTenantAdmin_False(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSONEnvelope(w, map[string]bool{"isAdmin": false}, nil)
	})
	isAdmin, err := cli.Users.IsTenantAdmin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if isAdmin {
		t.Error("expected isAdmin=false")
	}
}

func TestUsersClient_IsTenantAdmin_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Users.IsTenantAdmin(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_IsRSAuth_True(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/rs-auth-check" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, map[string]bool{"isAuth": true}, nil)
	})
	isAuth, err := cli.Users.IsRSAuth(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !isAuth {
		t.Error("expected isAuth=true")
	}
}

func TestUsersClient_IsRSAuth_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Users.IsRSAuth(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_IsViewerAuth_True(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/viewer-auth-check" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, map[string]bool{"isAuth": true}, nil)
	})
	isAuth, err := cli.Users.IsViewerAuth(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !isAuth {
		t.Error("expected isAuth=true")
	}
}

func TestUsersClient_IsViewerAuth_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	_, err := cli.Users.IsViewerAuth(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- Onboarding ----

func TestUsersClient_OnboardingVerify_Success(t *testing.T) {
	resp := types.LoginResponse{Token: "onboard-session"}
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/onboarding/verify" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, resp, nil)
	})
	req := OnboardingVerifyRequest{Data: OnboardingVerifyData{Token: "onboard-token", Password: "pass"}}
	result, err := cli.Users.OnboardingVerify(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "onboard-session" {
		t.Errorf("unexpected token: %q", result.Token)
	}
}

func TestUsersClient_OnboardingVerify_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Users.OnboardingVerify(context.Background(), OnboardingVerifyRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_OnboardingValidateToken_Valid(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/onboarding/validate-token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("token") != "valid-token" {
			t.Errorf("expected token param")
		}
		writeJSONEnvelope(w, map[string]bool{"valid": true}, nil)
	})
	valid, err := cli.Users.OnboardingValidateToken(context.Background(), "valid-token")
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Error("expected valid=true")
	}
}

func TestUsersClient_OnboardingValidateToken_Invalid(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSONEnvelope(w, map[string]bool{"valid": false}, nil)
	})
	valid, err := cli.Users.OnboardingValidateToken(context.Background(), "expired-token")
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Error("expected valid=false")
	}
}

func TestUsersClient_OnboardingValidateToken_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusBadRequest, nil)
	})
	_, err := cli.Users.OnboardingValidateToken(context.Background(), "bad-token")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUsersClient_OnboardingSendVerificationEmail_Success(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/web/api/v2.1/users/onboarding/send-verification-email" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSONEnvelope(w, nil, nil)
	})
	req := SendVerificationEmailRequest{Filter: BulkUsersFilter{Email: "alice@example.com"}}
	err := cli.Users.OnboardingSendVerificationEmail(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsersClient_OnboardingSendVerificationEmail_Error(t *testing.T) {
	_, cli := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeErrorEnvelope(w, http.StatusForbidden, nil)
	})
	err := cli.Users.OnboardingSendVerificationEmail(context.Background(), SendVerificationEmailRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// ---- ListUsersParams.values() ----

func TestListUsersParams_Values_Empty(t *testing.T) {
	p := &ListUsersParams{}
	vals := p.values()
	if len(vals) != 0 {
		t.Errorf("expected empty values, got %v", vals)
	}
}

func TestListUsersParams_Values_AllFields(t *testing.T) {
	p := &ListUsersParams{
		ListParams:               ListParams{Limit: IntPtr(50)},
		SiteIDs:                  []string{"s1"},
		AccountIDs:               []string{"acc1"},
		IDs:                      []string{"u1"},
		RoleIDs:                  []string{"r1"},
		Source:                   "mgmt",
		Sources:                  []string{"mgmt", "sso_saml"},
		Email:                    "alice@example.com",
		EmailContains:            []string{"example"},
		EmailReadOnly:            BoolPtr(false),
		FullName:                 "Alice Smith",
		FullNameContains:         []string{"Alice"},
		FullNameReadOnly:         BoolPtr(true),
		TwoFAEnabled:             BoolPtr(true),
		TwoFAStatus:              "enabled",
		TwoFAStatuses:            []string{"enabled", "disabled"},
		PrimaryTwoFAMethod:       "totp",
		EmailVerified:            BoolPtr(true),
		CanGenerateAPIToken:      BoolPtr(true),
		HasValidAPIToken:         BoolPtr(false),
		Query:                    "search",
		GroupsReadOnly:           BoolPtr(false),
		FirstLogin:               "2024-01-01",
		LastLogin:                "2024-06-01",
		DateJoined:               "2023-01-01",
		LastActivationLt:         "2024-07-01",
		LastActivationLte:        "2024-07-02",
		LastActivationGt:         "2024-05-01",
		LastActivationGte:        "2024-05-02",
		LastActivationBetween:    "2024-05-01_2024-07-01",
		APITokenExpiresAtLt:      "2025-01-01",
		APITokenExpiresAtLte:     "2025-01-02",
		APITokenExpiresAtGt:      "2024-12-01",
		APITokenExpiresAtGte:     "2024-12-02",
		APITokenExpiresAtBetween: "2024-12-01_2025-01-01",
		CreatedAtLt:              "2024-02-01",
		CreatedAtLte:             "2024-02-02",
		CreatedAtGt:              "2023-12-01",
		CreatedAtGte:             "2023-12-02",
		CreatedAtBetween:         "2023-12-01_2024-02-01",
	}
	vals := p.values()

	checks := map[string]string{
		"limit":                        "50",
		"siteIds":                      "s1",
		"accountIds":                   "acc1",
		"ids":                          "u1",
		"roleIds":                      "r1",
		"source":                       "mgmt",
		"sources":                      "mgmt,sso_saml",
		"email":                        "alice@example.com",
		"email__contains":              "example",
		"emailReadOnly":                "false",
		"fullName":                     "Alice Smith",
		"fullName__contains":           "Alice",
		"fullNameReadOnly":             "true",
		"twoFaEnabled":                 "true",
		"twoFaStatus":                  "enabled",
		"twoFaStatuses":                "enabled,disabled",
		"primaryTwoFaMethod":           "totp",
		"emailVerified":                "true",
		"canGenerateApiToken":          "true",
		"hasValidApiToken":             "false",
		"query":                        "search",
		"groupsReadOnly":               "false",
		"firstLogin":                   "2024-01-01",
		"lastLogin":                    "2024-06-01",
		"dateJoined":                   "2023-01-01",
		"lastActivation__lt":           "2024-07-01",
		"lastActivation__lte":          "2024-07-02",
		"lastActivation__gt":           "2024-05-01",
		"lastActivation__gte":          "2024-05-02",
		"lastActivation__between":      "2024-05-01_2024-07-01",
		"apiTokenExpiresAt__lt":        "2025-01-01",
		"apiTokenExpiresAt__lte":       "2025-01-02",
		"apiTokenExpiresAt__gt":        "2024-12-01",
		"apiTokenExpiresAt__gte":       "2024-12-02",
		"apiTokenExpiresAt__between":   "2024-12-01_2025-01-01",
		"createdAt__lt":                "2024-02-01",
		"createdAt__lte":               "2024-02-02",
		"createdAt__gt":                "2023-12-01",
		"createdAt__gte":               "2023-12-02",
		"createdAt__between":           "2023-12-01_2024-02-01",
	}
	for key, want := range checks {
		if got := vals.Get(key); got != want {
			t.Errorf("param %q: got %q, want %q", key, got, want)
		}
	}
}
