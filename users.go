package sentinelone

import (
	"context"
	"fmt"
	"net/url"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// UsersClient provides access to the Users API group.
// Access it via [Client.Users].
type UsersClient struct{ c *Client }

// -- User CRUD --

// List returns a paginated list of users visible to the authenticated user,
// filtered by the optional params.
//
// API: GET /web/api/v2.1/users
// Required permission: Users.view
//
// Pass nil for params to use the API defaults (limit 10, no filters).
func (u *UsersClient) List(ctx context.Context, params *ListUsersParams) ([]types.User, *types.Pagination, error) {
	var paramVals url.Values
	if params != nil {
		paramVals = params.values()
	}

	var users []types.User

	pag, err := u.c.get(ctx, "/users", paramVals, &users)
	if err != nil {
		return nil, nil, err
	}

	return users, pag, nil
}

// Get returns the user with the given userID.
//
// API: GET /web/api/v2.1/users/{user_id}
// Required permission: Users.view
func (u *UsersClient) Get(ctx context.Context, userID string) (*types.User, error) {
	var user types.User

	_, err := u.c.get(ctx, "/users/"+userID, nil, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create creates a new user account.
//
// API: POST /web/api/v2.1/users
// Required permission: Users.create
//
// Email, FullName, and Scope are required.  Scope must be one of "site",
// "account", or "tenant".  Assign roles via ScopeRoles (preferred) or the
// deprecated SiteRoles / TenantRoles fields.  If automatic onboarding is
// enabled for the tenant, omit Password — the user will receive an invitation
// email instead.
func (u *UsersClient) Create(ctx context.Context, req CreateUserRequest) (*types.User, error) {
	var user types.User

	_, err := u.c.post(ctx, "/users", req, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates mutable fields on an existing user.
//
// API: PUT /web/api/v2.1/users/{user_id}
// Required permission: Users.update
//
// Scope is required in req.Data.  Only non-zero fields are applied.  To
// change a password, set both CurrentPassword and Password.  Role assignments
// are replaced in full when ScopeRoles is non-empty.
func (u *UsersClient) Update(ctx context.Context, userID string, req UpdateUserRequest) (*types.User, error) {
	var user types.User

	_, err := u.c.put(ctx, "/users/"+userID, req, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Delete permanently deletes the user with the given userID.
//
// API: DELETE /web/api/v2.1/users/{user_id}
// Required permission: Users.delete
func (u *UsersClient) Delete(ctx context.Context, userID string) error {
	_, err := u.c.delete(ctx, "/users/"+userID, nil, nil)

	return err
}

// BulkDelete permanently deletes all users that match the filter in req.
// Use with caution — this operation is irreversible.  Narrow the target set
// with req.Filter.IDs, req.Filter.Email, or scope-based criteria before
// calling.
//
// API: POST /web/api/v2.1/users/delete-users
// Required permission: Users.delete
func (u *UsersClient) BulkDelete(ctx context.Context, req BulkUsersActionRequest) error {
	_, err := u.c.post(ctx, "/users/delete-users", req, nil)

	return err
}

// -- API token management --

// GetAPITokenDetails returns the creation and expiry timestamps of the API
// token belonging to the specified user.  It does not return the token value.
//
// API: GET /web/api/v2.1/users/{user_id}/api-token-details
// Required permission: Users.view
func (u *UsersClient) GetAPITokenDetails(ctx context.Context, userID string) (*types.APITokenDetail, error) {
	var detail types.APITokenDetail

	_, err := u.c.get(ctx, fmt.Sprintf("/users/%s/api-token-details", userID), nil, &detail)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

// GetAPITokenDetailsByToken looks up the creation and expiry timestamps for an
// API token supplied as a string.  Useful for validating a token before use.
//
// API: POST /web/api/v2.1/users/api-token-details
// Required permission: Users.view
func (u *UsersClient) GetAPITokenDetailsByToken(
	ctx context.Context,
	req GetAPITokenDetailsRequest,
) (*types.APITokenDetail, error) {
	var detail types.APITokenDetail

	_, err := u.c.post(ctx, "/users/api-token-details", req, &detail)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

// GenerateAPIToken generates a new API token for the currently authenticated
// user, replacing any existing token.  The plaintext token is returned in the
// response and cannot be retrieved again; store it securely.
//
// API: POST /web/api/v2.1/users/generate-api-token
// Required permission: (self — no special permission beyond being logged in)
func (u *UsersClient) GenerateAPIToken(
	ctx context.Context,
	req GenerateAPITokenRequest,
) (*types.APITokenResponse, error) {
	var token types.APITokenResponse

	_, err := u.c.post(ctx, "/users/generate-api-token", req, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// RevokeAPIToken invalidates the API token for the user identified by
// req.Data.UserID.  The user will need to generate a new token before making
// API calls again.
//
// API: POST /web/api/v2.1/users/revoke-api-token
// Required permission: Users.revokeApiToken
func (u *UsersClient) RevokeAPIToken(ctx context.Context, req UserIDRequest) error {
	_, err := u.c.post(ctx, "/users/revoke-api-token", req, nil)

	return err
}

// -- 2FA management --

// Enable2FA enables the two-factor authentication requirement for the user
// identified by req.Data.UserID.  The user must subsequently complete
// enrollment via [UsersClient.Enroll2FA].
//
// API: POST /web/api/v2.1/users/2fa/enable
// Required permission: Users.twoFa
func (u *UsersClient) Enable2FA(ctx context.Context, req UserIDRequest) error {
	_, err := u.c.post(ctx, "/users/2fa/enable", req, nil)

	return err
}

// Disable2FA removes the two-factor authentication requirement for the user
// identified by req.Data.UserID.
//
// API: POST /web/api/v2.1/users/2fa/disable
// Required permission: Users.twoFa
func (u *UsersClient) Disable2FA(ctx context.Context, req UserIDRequest) error {
	_, err := u.c.post(ctx, "/users/2fa/disable", req, nil)

	return err
}

// Enroll2FA initiates the 2FA enrollment flow for one or more users.  The
// response contains a TOTP secret and QR code URL that the user should scan
// with an authenticator app.
//
// API: POST /web/api/v2.1/users/enroll-2fa
// Required permission: Users.twoFa
func (u *UsersClient) Enroll2FA(ctx context.Context, req UserIDsRequest) (*types.EnrollTFAResponse, error) {
	var resp types.EnrollTFAResponse

	_, err := u.c.post(ctx, "/users/enroll-2fa", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Reset2FA clears a user's enrolled 2FA device, forcing them to re-enroll on
// their next login.  Use this when a user loses access to their authenticator.
//
// API: POST /web/api/v2.1/users/reset-2fa
// Required permission: Users.twoFa
func (u *UsersClient) Reset2FA(ctx context.Context, req ResetTFARequest) error {
	_, err := u.c.post(ctx, "/users/reset-2fa", req, nil)

	return err
}

// Delete2FA removes a user's 2FA configuration entirely, including any enrolled
// device.  After this call, 2FA is no longer configured for the user.
//
// API: POST /web/api/v2.1/users/delete-2fa
// Required permission: Users.twoFa
func (u *UsersClient) Delete2FA(ctx context.Context, req DeleteTFARequest) error {
	_, err := u.c.post(ctx, "/users/delete-2fa", req, nil)

	return err
}

// -- Password management --

// ChangePassword changes the password for the currently authenticated user.
// CurrentPassword is required for verification; if 2FA is enabled, TwoFACode
// must also be provided.
//
// API: POST /web/api/v2.1/users/change-password
// Required permission: (self — no special permission beyond being logged in)
func (u *UsersClient) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	_, err := u.c.post(ctx, "/users/change-password", req, nil)

	return err
}

// SendResetPasswordEmail triggers a password-reset email for all users that
// match req.Filter.  Typically used with req.Filter.Email or req.Filter.IDs
// to target a single user.
//
// API: POST /web/api/v2.1/users/login/send-reset-password-email
// Required permission: Users.resetPassword
func (u *UsersClient) SendResetPasswordEmail(ctx context.Context, req SendResetPasswordRequest) error {
	_, err := u.c.post(ctx, "/users/login/send-reset-password-email", req, nil)

	return err
}

// ForceResetPasswordOnLogin flags all users matching req.Filter so that they
// are required to choose a new password at their next login attempt.
//
// API: POST /web/api/v2.1/users/login/force-reset-password-on-login
// Required permission: Users.resetPassword
func (u *UsersClient) ForceResetPasswordOnLogin(ctx context.Context, req ForceResetPasswordRequest) error {
	_, err := u.c.post(ctx, "/users/login/force-reset-password-on-login", req, nil)

	return err
}

// -- App / iframe integration --

// EnableApp enables or disables an integrated application for users.
// Set req.Data.Enable to true to enable the app, false to disable it.
//
// API: POST /web/api/v2.1/users/enable-app
// Required permission: Users.enableApp
func (u *UsersClient) EnableApp(ctx context.Context, req EnableAppRequest) error {
	_, err := u.c.post(ctx, "/users/enable-app", req, nil)

	return err
}

// RequestApp submits a request for the authenticated user to access an
// application.  The response URL redirects to the app's access-request flow.
//
// API: POST /web/api/v2.1/users/request-app
// Required permission: (self — no special permission beyond being logged in)
func (u *UsersClient) RequestApp(ctx context.Context, req RequestAppRequest) (*types.RequestAppResponse, error) {
	var resp types.RequestAppResponse

	_, err := u.c.post(ctx, "/users/request-app", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GenerateIFrameToken generates a short-lived token that allows embedding a
// SentinelOne console view in an iframe.  Scope the token to an account or
// site by setting req.Data.AccountID or req.Data.SiteID.
//
// API: POST /web/api/v2.1/users/generate-iframe-token
// Required permission: Users.generateIframeToken
func (u *UsersClient) GenerateIFrameToken(
	ctx context.Context,
	req IFrameUserRequest,
) (*types.IFrameTokenResponse, error) {
	var token types.IFrameTokenResponse

	_, err := u.c.post(ctx, "/users/generate-iframe-token", req, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// -- Authentication --

// Login authenticates with a username (email) and password, returning a session
// token on success.  If the account requires 2FA, the response Status will be
// "2fa_required" and TwoFAMethod will indicate which method to use; pass the
// resulting token and a TOTP code to [UsersClient.LoginContinue].
//
// The returned token can be used to construct an authenticated [Client] with
// [NewClient].
//
// API: POST /web/api/v2.1/users/login
func (u *UsersClient) Login(ctx context.Context, req LoginRequest) (*types.LoginResponse, error) {
	var resp types.LoginResponse

	_, err := u.c.post(ctx, "/users/login", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// LoginContinue completes a multi-step login by providing the second factor.
// Populate req.Data.Token with the token from [UsersClient.Login] and
// req.Data.Code with the user's TOTP or SMS code.
//
// API: POST /web/api/v2.1/users/login-continue
func (u *UsersClient) LoginContinue(
	ctx context.Context,
	req LoginContinueRequest,
) (*types.LoginContinueResponse, error) {
	var resp types.LoginContinueResponse

	_, err := u.c.post(ctx, "/users/login-continue", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// LoginByToken authenticates using a short-lived login token (e.g. one
// received in a password-reset or invitation email).
//
// API: GET /web/api/v2.1/users/login/by-token
func (u *UsersClient) LoginByToken(ctx context.Context, token string) (*types.LoginResponse, error) {
	queryParams := url.Values{}
	queryParams.Set("token", token)

	var resp types.LoginResponse

	_, err := u.c.get(ctx, "/users/login/by-token", queryParams, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// LoginByAPIToken authenticates using a long-lived API token, returning a
// short-lived session token suitable for UI-style calls.  For pure API usage
// it is simpler to pass the API token directly to [NewClient].
//
// API: POST /web/api/v2.1/users/login/by-api-token
func (u *UsersClient) LoginByAPIToken(ctx context.Context, req LoginByAPITokenRequest) (*types.LoginResponse, error) {
	var resp types.LoginResponse

	_, err := u.c.post(ctx, "/users/login/by-api-token", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// LoginSSO initiates a SAML2 SSO login flow and returns the redirect URL that
// the user's browser should be sent to in order to complete authentication
// with the identity provider.
//
// API: GET /web/api/v2.1/users/login/sso-saml2
func (u *UsersClient) LoginSSO(ctx context.Context) (string, error) {
	var result struct {
		RedirectURL string `json:"redirectUrl,omitempty"`
	}

	_, err := u.c.get(ctx, "/users/login/sso-saml2", nil, &result)
	if err != nil {
		return "", err
	}

	return result.RedirectURL, nil
}

// LoginSSOForScope initiates a SAML2 SSO login scoped to a specific account or
// site (identified by scopeID), returning the identity-provider redirect URL.
//
// API: POST /web/api/v2.1/users/login/sso-saml2/{scope_id}
func (u *UsersClient) LoginSSOForScope(ctx context.Context, scopeID string) (string, error) {
	var result struct {
		RedirectURL string `json:"redirectUrl,omitempty"`
	}

	_, err := u.c.post(ctx, "/users/login/sso-saml2/"+scopeID, struct{}{}, &result)
	if err != nil {
		return "", err
	}

	return result.RedirectURL, nil
}

// SetPassword sets a new password for a user using a password-reset token.
// Token comes from a password-reset email link; the response contains a
// session token if the operation succeeds.
//
// API: POST /web/api/v2.1/users/login/set-password
func (u *UsersClient) SetPassword(ctx context.Context, req SetPasswordRequest) (*types.SetPasswordResponse, error) {
	var resp types.SetPasswordResponse

	_, err := u.c.post(ctx, "/users/login/set-password", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Logout invalidates the current session token.  After this call, the client's
// token can no longer be used to authenticate requests.
//
// API: POST /web/api/v2.1/users/logout
func (u *UsersClient) Logout(ctx context.Context) error {
	_, err := u.c.post(ctx, "/users/logout", struct{}{}, nil)

	return err
}

// -- Session management --

// ElevateSession re-verifies the caller's identity by requiring their password,
// upgrading the current session token to an elevated-privilege token.  Use
// this before sensitive operations that require step-up authentication.
//
// API: POST /web/api/v2.1/users/auth/elevate
// Required permission: (self — no special permission beyond being logged in)
func (u *UsersClient) ElevateSession(
	ctx context.Context,
	req ElevateSessionRequest,
) (*types.ElevateSessionResponse, error) {
	var resp types.ElevateSessionResponse

	_, err := u.c.post(ctx, "/users/auth/elevate", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// SSOReAuth initiates a SAML2 re-authentication challenge for the current
// session, typically triggered when an elevated-privilege action is required.
// Returns the identity-provider redirect URL.
//
// API: GET /web/api/v2.1/users/sso-saml2/re-auth
func (u *UsersClient) SSOReAuth(ctx context.Context) (string, error) {
	var result struct {
		RedirectURL string `json:"redirectUrl,omitempty"`
	}

	_, err := u.c.get(ctx, "/users/sso-saml2/re-auth", nil, &result)
	if err != nil {
		return "", err
	}

	return result.RedirectURL, nil
}

// AuthEULA records that the authenticated user has accepted the End User
// License Agreement.  This is typically called once after first login when
// the server responds with a status indicating EULA acceptance is required.
//
// API: POST /web/api/v2.1/users/auth/eula
func (u *UsersClient) AuthEULA(ctx context.Context) error {
	_, err := u.c.post(ctx, "/users/auth/eula", struct{}{}, nil)

	return err
}

// AuthApp completes an OAuth-style app authorization flow by exchanging an
// authorization code for a session token.
//
// API: POST /web/api/v2.1/users/auth/app
func (u *UsersClient) AuthApp(ctx context.Context, req AuthAppRequest) (*types.LoginResponse, error) {
	var resp types.LoginResponse

	_, err := u.c.post(ctx, "/users/auth/app", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// IsTenantAdmin reports whether the currently authenticated user holds a
// tenant-admin role.  Returns true only when that role is confirmed.
//
// API: GET /web/api/v2.1/users/tenant-admin-auth-check
func (u *UsersClient) IsTenantAdmin(ctx context.Context) (bool, error) {
	var result struct {
		IsAdmin bool `json:"isAdmin"`
	}

	_, err := u.c.get(ctx, "/users/tenant-admin-auth-check", nil, &result)
	if err != nil {
		return false, err
	}

	return result.IsAdmin, nil
}

// IsRSAuth reports whether the currently authenticated user has Remote Shell
// (RS) authorization.
//
// API: GET /web/api/v2.1/users/rs-auth-check
func (u *UsersClient) IsRSAuth(ctx context.Context) (bool, error) {
	var result struct {
		IsAuth bool `json:"isAuth"`
	}

	_, err := u.c.get(ctx, "/users/rs-auth-check", nil, &result)
	if err != nil {
		return false, err
	}

	return result.IsAuth, nil
}

// IsViewerAuth reports whether the currently authenticated user holds at
// least a viewer-level role.
//
// API: GET /web/api/v2.1/users/viewer-auth-check
func (u *UsersClient) IsViewerAuth(ctx context.Context) (bool, error) {
	var result struct {
		IsAuth bool `json:"isAuth"`
	}

	_, err := u.c.get(ctx, "/users/viewer-auth-check", nil, &result)
	if err != nil {
		return false, err
	}

	return result.IsAuth, nil
}

// -- Onboarding --

// OnboardingVerify validates the onboarding token sent to a new user's email
// and optionally sets their initial password.  On success the response
// contains a session token.
//
// API: POST /web/api/v2.1/users/onboarding/verify
func (u *UsersClient) OnboardingVerify(ctx context.Context, req OnboardingVerifyRequest) (*types.LoginResponse, error) {
	var resp types.LoginResponse

	_, err := u.c.post(ctx, "/users/onboarding/verify", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// OnboardingValidateToken checks whether the given onboarding token is still
// valid (not expired, not already used).  Returns true if valid.
//
// API: GET /web/api/v2.1/users/onboarding/validate-token
func (u *UsersClient) OnboardingValidateToken(ctx context.Context, token string) (bool, error) {
	queryParams := url.Values{}
	queryParams.Set("token", token)

	var result struct {
		Valid bool `json:"valid,omitempty"`
	}

	_, err := u.c.get(ctx, "/users/onboarding/validate-token", queryParams, &result)
	if err != nil {
		return false, err
	}

	return result.Valid, nil
}

// OnboardingSendVerificationEmail resends the verification/onboarding email to
// all users that match req.Filter.  Useful when a user did not receive or has
// lost their original invitation email.
//
// API: POST /web/api/v2.1/users/onboarding/send-verification-email
// Required permission: Users.create
func (u *UsersClient) OnboardingSendVerificationEmail(ctx context.Context, req SendVerificationEmailRequest) error {
	_, err := u.c.post(ctx, "/users/onboarding/send-verification-email", req, nil)

	return err
}
