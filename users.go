package sentinelone

import (
	"context"
	"fmt"
	"net/url"
)

// -- User types --

// UserAPIToken holds API token metadata for a user.
type UserAPIToken struct {
	CreatedAt string `json:"createdAt,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// UserScopeRole links a scope (account/site/group) to a role assignment.
type UserScopeRole struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	AccountName string   `json:"accountName"`
	Roles       []string `json:"roles,omitempty"`
	RoleName    string   `json:"roleName,omitempty"`
	RoleID      string   `json:"roleId,omitempty"`
}

// UserSiteRole is a deprecated site-scoped role assignment.
type UserSiteRole struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles,omitempty"`
	RoleName string   `json:"roleName,omitempty"`
	RoleID   string   `json:"roleId,omitempty"`
}

// User represents a SentinelOne management user.
type User struct {
	ID                    string          `json:"id,omitempty"`
	Source                string          `json:"source,omitempty"`
	GlobalUserID          *string         `json:"globalUserId,omitempty"`
	GlobalOrganizationID  *string         `json:"globalOrganizationId,omitempty"`
	Email                 *string         `json:"email,omitempty"`
	EmailReadOnly         bool            `json:"emailReadOnly,omitempty"`
	FullName              *string         `json:"fullName,omitempty"`
	FullNameReadOnly      bool            `json:"fullNameReadOnly,omitempty"`
	FirstLogin            string          `json:"firstLogin,omitempty"`
	LastLogin             string          `json:"lastLogin,omitempty"`
	DateJoined            string          `json:"dateJoined,omitempty"`
	TwoFAEnabled          bool            `json:"twoFaEnabled,omitempty"`
	TwoFAEnabledReadOnly  bool            `json:"twoFaEnabledReadOnly,omitempty"`
	TwoFAConfigured       bool            `json:"twoFaConfigured,omitempty"`
	TwoFAStatus           string          `json:"twoFaStatus,omitempty"`
	PrimaryTwoFAMethod    string          `json:"primaryTwoFaMethod,omitempty"`
	EmailVerified         bool            `json:"emailVerified,omitempty"`
	Scope                 string          `json:"scope"`
	APIToken              *UserAPIToken   `json:"apiToken,omitempty"`
	AgreedEULA            bool            `json:"agreedEula,omitempty"`
	AgreementURL          string          `json:"agreementUrl,omitempty"`
	CanGenerateAPIToken   bool            `json:"canGenerateApiToken,omitempty"`
	IsSystem              bool            `json:"isSystem,omitempty"`
	ScopeRoles            []UserScopeRole `json:"scopeRoles,omitempty"`
	SiteRoles             []UserSiteRole  `json:"siteRoles,omitempty"`
	TenantRoles           []string        `json:"tenantRoles,omitempty"`

	// Deprecated
	LowestRole     string `json:"lowestRole,omitempty"`
	GroupsReadOnly bool   `json:"groupsReadOnly,omitempty"`
	AllowRemoteShell bool  `json:"allowRemoteShell,omitempty"`
}

// LoginResponse holds the result of a successful login.
type LoginResponse struct {
	Token       string `json:"token"`
	Status      string `json:"status,omitempty"`
	TwoFAMethod string `json:"twoFaMethod,omitempty"`
	CSRF        string `json:"csrf,omitempty"`
}

// LoginContinueResponse holds the result of a login-continue operation (e.g., after 2FA).
type LoginContinueResponse struct {
	Token  string `json:"token,omitempty"`
	Status string `json:"status,omitempty"`
}

// APITokenResponse holds a user API token.
type APITokenResponse struct {
	Token string `json:"token"`
}

// APITokenDetail holds metadata about a user's API token.
type APITokenDetail struct {
	CreatedAt string `json:"createdAt,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// EnrollTFAResponse holds the 2FA enrollment data (QR code / secret).
type EnrollTFAResponse struct {
	Secret string `json:"secret,omitempty"`
	QRCode string `json:"qrCode,omitempty"`
}

// IFrameTokenResponse holds a generated iframe access token.
type IFrameTokenResponse struct {
	Token string `json:"token,omitempty"`
}

// ElevateSessionResponse holds the result of an elevated-session request.
type ElevateSessionResponse struct {
	Token string `json:"token,omitempty"`
}

// RequestAppResponse holds the app-access request result.
type RequestAppResponse struct {
	URL string `json:"url,omitempty"`
}

// SetPasswordResponse holds the result of a set-password call.
type SetPasswordResponse struct {
	Token string `json:"token,omitempty"`
}

// -- Request types --

// CreateUserRequest is the request body for POST /users.
type CreateUserRequest struct {
	Data CreateUserData `json:"data"`
}

// CreateUserData holds the fields for creating a user.
type CreateUserData struct {
	Email              string          `json:"email"`
	FullName           string          `json:"fullName"`
	Scope              string          `json:"scope"`
	Password           string          `json:"password,omitempty"`
	TwoFAEnabled       *bool           `json:"twoFaEnabled,omitempty"`
	ForceCredentialAuth *bool          `json:"forceCredentialAuth,omitempty"`
	ScopeRoles         []UserScopeRole `json:"scopeRoles,omitempty"`
	SiteRoles          []UserSiteRole  `json:"siteRoles,omitempty"`
	TenantRoles        []string        `json:"tenantRoles,omitempty"`
}

// UpdateUserRequest is the request body for PUT /users/{id}.
type UpdateUserRequest struct {
	Data UpdateUserData `json:"data"`
}

// UpdateUserData holds the fields for updating a user.
type UpdateUserData struct {
	Scope               string          `json:"scope"`
	ID                  string          `json:"id,omitempty"`
	Email               string          `json:"email,omitempty"`
	FullName            string          `json:"fullName,omitempty"`
	Password            string          `json:"password,omitempty"`
	CurrentPassword     string          `json:"currentPassword,omitempty"`
	TwoFAEnabled        *bool           `json:"twoFaEnabled,omitempty"`
	TwoFACode           string          `json:"twoFaCode,omitempty"`
	CanGenerateAPIToken *bool           `json:"canGenerateApiToken,omitempty"`
	ScopeRoles          []UserScopeRole `json:"scopeRoles,omitempty"`
	SiteRoles           []UserSiteRole  `json:"siteRoles,omitempty"`
	TenantRoles         []string        `json:"tenantRoles,omitempty"`
}

// LoginRequest is the request body for POST /users/login.
type LoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe *bool  `json:"rememberMe,omitempty"`
}

// LoginContinueRequest is the request body for POST /users/login-continue.
type LoginContinueRequest struct {
	Data LoginContinueData `json:"data"`
}

// LoginContinueData holds 2FA or continuation data.
type LoginContinueData struct {
	Token  string `json:"token,omitempty"`
	Code   string `json:"code,omitempty"`
	Method string `json:"method,omitempty"`
}

// LoginByAPITokenRequest is the request body for POST /users/login/by-api-token.
type LoginByAPITokenRequest struct {
	Data LoginByAPITokenData `json:"data"`
}

// LoginByAPITokenData holds the API token for token-based login.
type LoginByAPITokenData struct {
	APIToken string `json:"apiToken"`
}

// GenerateAPITokenRequest is the request body for POST /users/generate-api-token.
type GenerateAPITokenRequest struct {
	Data GenerateAPITokenData `json:"data"`
}

// GenerateAPITokenData holds options for token generation.
type GenerateAPITokenData struct {
	ForceLegacy *bool `json:"forceLegacy,omitempty"`
}

// GetAPITokenDetailsRequest is the request body for POST /users/api-token-details.
type GetAPITokenDetailsRequest struct {
	Data GetAPITokenDetailsData `json:"data"`
}

// GetAPITokenDetailsData holds the token to look up.
type GetAPITokenDetailsData struct {
	Token string `json:"token"`
}

// ChangePasswordRequest is the request body for POST /users/change-password.
type ChangePasswordRequest struct {
	Data ChangePasswordData `json:"data"`
}

// ChangePasswordData holds the password change fields.
type ChangePasswordData struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
	TwoFACode       string `json:"twoFaCode,omitempty"`
}

// SetPasswordRequest is the request body for POST /users/login/set-password.
type SetPasswordRequest struct {
	Data SetPasswordData `json:"data"`
}

// SetPasswordData holds the new password and reset token.
type SetPasswordData struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// UserIDRequest is used for single-user operations like enable/disable 2FA.
type UserIDRequest struct {
	Data UserIDData `json:"data"`
}

// UserIDData holds a user ID.
type UserIDData struct {
	UserID string `json:"userId"`
}

// UserIDsRequest is used for multi-user operations like enroll-2fa.
type UserIDsRequest struct {
	Data UserIDsData `json:"data"`
}

// UserIDsData holds a list of user IDs.
type UserIDsData struct {
	UserIDs []string `json:"userIds"`
}

// BulkUsersFilter specifies which users to target for bulk operations.
type BulkUsersFilter struct {
	IDs                    []string `json:"ids,omitempty"`
	Email                  string   `json:"email,omitempty"`
	EmailContains          []string `json:"email__contains,omitempty"`
	FullName               string   `json:"fullName,omitempty"`
	FullNameContains       []string `json:"fullName__contains,omitempty"`
	Source                 string   `json:"source,omitempty"`
	Sources                []string `json:"sources,omitempty"`
	TwoFAEnabled           *bool    `json:"twoFaEnabled,omitempty"`
	TwoFAStatus            string   `json:"twoFaStatus,omitempty"`
	TwoFAStatuses          []string `json:"twoFaStatuses,omitempty"`
	EmailVerified          *bool    `json:"emailVerified,omitempty"`
	EmailReadOnly          *bool    `json:"emailReadOnly,omitempty"`
	FullNameReadOnly       *bool    `json:"fullNameReadOnly,omitempty"`
	GroupsReadOnly         *bool    `json:"groupsReadOnly,omitempty"`
	RoleIDs                []string `json:"roleIds,omitempty"`
	CanGenerateAPIToken    *bool    `json:"canGenerateApiToken,omitempty"`
	HasValidAPIToken       *bool    `json:"hasValidApiToken,omitempty"`
	Query                  string   `json:"query,omitempty"`
	LastActivationLt       string   `json:"lastActivation__lt,omitempty"`
	LastActivationLte      string   `json:"lastActivation__lte,omitempty"`
	LastActivationGt       string   `json:"lastActivation__gt,omitempty"`
	LastActivationGte      string   `json:"lastActivation__gte,omitempty"`
	LastActivationBetween  string   `json:"lastActivation__between,omitempty"`
	APITokenExpiresAtLt    string   `json:"apiTokenExpiresAt__lt,omitempty"`
	APITokenExpiresAtLte   string   `json:"apiTokenExpiresAt__lte,omitempty"`
	APITokenExpiresAtGt    string   `json:"apiTokenExpiresAt__gt,omitempty"`
	APITokenExpiresAtGte   string   `json:"apiTokenExpiresAt__gte,omitempty"`
	APITokenExpiresAtBetween string `json:"apiTokenExpiresAt__between,omitempty"`
	CreatedAtLt            string   `json:"createdAt__lt,omitempty"`
	CreatedAtLte           string   `json:"createdAt__lte,omitempty"`
	CreatedAtGt            string   `json:"createdAt__gt,omitempty"`
	CreatedAtGte           string   `json:"createdAt__gte,omitempty"`
	CreatedAtBetween       string   `json:"createdAt__between,omitempty"`
}

// BulkUsersActionRequest is used for bulk user operations like delete-users.
type BulkUsersActionRequest struct {
	Filter BulkUsersFilter `json:"filter"`
	Data   *struct{}       `json:"data,omitempty"`
}

// ElevateSessionRequest is the request body for POST /users/auth/elevate.
type ElevateSessionRequest struct {
	Data ElevateSessionData `json:"data"`
}

// ElevateSessionData holds the password for session elevation.
type ElevateSessionData struct {
	Password string `json:"password"`
}

// AuthAppRequest is the request body for POST /users/auth/app.
type AuthAppRequest struct {
	Data AuthAppData `json:"data"`
}

// AuthAppData holds the authorization code for app auth.
type AuthAppData struct {
	Code string `json:"code"`
}

// EnableAppRequest is the request body for POST /users/enable-app.
type EnableAppRequest struct {
	Data EnableAppData `json:"data"`
}

// EnableAppData holds app enable parameters.
type EnableAppData struct {
	AppID  string `json:"appId"`
	Enable bool   `json:"enable"`
}

// RequestAppRequest is the request body for POST /users/request-app.
type RequestAppRequest struct {
	CurrentPassword string `json:"currentPassword,omitempty"`
}

// ResetTFARequest is the request body for POST /users/reset-2fa.
type ResetTFARequest struct {
	Data ResetTFAData `json:"data"`
}

// ResetTFAData holds the user ID for 2FA reset.
type ResetTFAData struct {
	UserID string `json:"userId"`
}

// DeleteTFARequest is the request body for POST /users/delete-2fa.
type DeleteTFARequest struct {
	Data DeleteTFAData `json:"data"`
}

// DeleteTFAData holds the user ID for 2FA deletion.
type DeleteTFAData struct {
	UserID string `json:"userId"`
}

// OnboardingVerifyRequest is the request body for POST /users/onboarding/verify.
type OnboardingVerifyRequest struct {
	Data OnboardingVerifyData `json:"data"`
}

// OnboardingVerifyData holds onboarding verification fields.
type OnboardingVerifyData struct {
	Token    string `json:"token"`
	Password string `json:"password,omitempty"`
}

// SendResetPasswordRequest is the request body for POST /users/login/send-reset-password-email.
type SendResetPasswordRequest struct {
	Filter BulkUsersFilter `json:"filter"`
}

// ForceResetPasswordRequest is the request body for POST /users/login/force-reset-password-on-login.
type ForceResetPasswordRequest struct {
	Filter BulkUsersFilter `json:"filter"`
}

// SendVerificationEmailRequest is the request body for POST /users/onboarding/send-verification-email.
type SendVerificationEmailRequest struct {
	Filter BulkUsersFilter `json:"filter"`
}

// IFrameUserRequest is the request body for POST /users/generate-iframe-token.
type IFrameUserRequest struct {
	Data IFrameUserData `json:"data"`
}

// IFrameUserData holds the iframe user parameters.
type IFrameUserData struct {
	AccountID  string `json:"accountId,omitempty"`
	SiteID     string `json:"siteId,omitempty"`
	Expiration string `json:"expiration,omitempty"`
}

// -- Filter params --

// ListUsersParams contains query parameters for GET /web/api/v2.1/users.
// All fields are optional; zero values are omitted from the request.
//
//   - SiteIDs / AccountIDs: limit results to users in specific scopes.
//   - IDs / RoleIDs: filter to specific user or role IDs.
//   - Email / EmailContains: exact or partial email address match.
//   - FullName / FullNameContains: exact or partial full name match.
//   - Source / Sources: "mgmt", "sso_saml", "active_directory", or "global".
//   - TwoFAEnabled / TwoFAStatus: filter by 2FA state.
//   - HasValidAPIToken / CanGenerateAPIToken: filter by API token status.
//   - LastActivation* / APITokenExpiresAt* / CreatedAt*: timestamp range filters
//     using the __lt/__gt/__lte/__gte/__between suffixes accepted by the API.
type ListUsersParams struct {
	ListParams
	SiteIDs                  []string
	AccountIDs               []string
	IDs                      []string
	RoleIDs                  []string
	Source                   string
	Sources                  []string
	Email                    string
	EmailContains            []string
	EmailReadOnly            *bool
	FullName                 string
	FullNameContains         []string
	FullNameReadOnly         *bool
	TwoFAEnabled             *bool
	TwoFAStatus              string
	TwoFAStatuses            []string
	PrimaryTwoFAMethod       string
	EmailVerified            *bool
	CanGenerateAPIToken      *bool
	HasValidAPIToken         *bool
	Query                    string
	GroupsReadOnly           *bool
	FirstLogin               string
	LastLogin                string
	DateJoined               string
	LastActivationLt         string
	LastActivationLte        string
	LastActivationGt         string
	LastActivationGte        string
	LastActivationBetween    string
	APITokenExpiresAtLt      string
	APITokenExpiresAtLte     string
	APITokenExpiresAtGt      string
	APITokenExpiresAtGte     string
	APITokenExpiresAtBetween string
	CreatedAtLt              string
	CreatedAtLte             string
	CreatedAtGt              string
	CreatedAtGte             string
	CreatedAtBetween         string
}

func (p *ListUsersParams) values() url.Values {
	v := p.ListParams.values()
	setStringSlice(v, "siteIds", p.SiteIDs)
	setStringSlice(v, "accountIds", p.AccountIDs)
	setStringSlice(v, "ids", p.IDs)
	setStringSlice(v, "roleIds", p.RoleIDs)
	setString(v, "source", &p.Source)
	setStringSlice(v, "sources", p.Sources)
	setString(v, "email", &p.Email)
	setStringSlice(v, "email__contains", p.EmailContains)
	setBool(v, "emailReadOnly", p.EmailReadOnly)
	setString(v, "fullName", &p.FullName)
	setStringSlice(v, "fullName__contains", p.FullNameContains)
	setBool(v, "fullNameReadOnly", p.FullNameReadOnly)
	setBool(v, "twoFaEnabled", p.TwoFAEnabled)
	setString(v, "twoFaStatus", &p.TwoFAStatus)
	setStringSlice(v, "twoFaStatuses", p.TwoFAStatuses)
	setString(v, "primaryTwoFaMethod", &p.PrimaryTwoFAMethod)
	setBool(v, "emailVerified", p.EmailVerified)
	setBool(v, "canGenerateApiToken", p.CanGenerateAPIToken)
	setBool(v, "hasValidApiToken", p.HasValidAPIToken)
	setString(v, "query", &p.Query)
	setBool(v, "groupsReadOnly", p.GroupsReadOnly)
	setString(v, "firstLogin", &p.FirstLogin)
	setString(v, "lastLogin", &p.LastLogin)
	setString(v, "dateJoined", &p.DateJoined)
	setString(v, "lastActivation__lt", &p.LastActivationLt)
	setString(v, "lastActivation__lte", &p.LastActivationLte)
	setString(v, "lastActivation__gt", &p.LastActivationGt)
	setString(v, "lastActivation__gte", &p.LastActivationGte)
	setString(v, "lastActivation__between", &p.LastActivationBetween)
	setString(v, "apiTokenExpiresAt__lt", &p.APITokenExpiresAtLt)
	setString(v, "apiTokenExpiresAt__lte", &p.APITokenExpiresAtLte)
	setString(v, "apiTokenExpiresAt__gt", &p.APITokenExpiresAtGt)
	setString(v, "apiTokenExpiresAt__gte", &p.APITokenExpiresAtGte)
	setString(v, "apiTokenExpiresAt__between", &p.APITokenExpiresAtBetween)
	setString(v, "createdAt__lt", &p.CreatedAtLt)
	setString(v, "createdAt__lte", &p.CreatedAtLte)
	setString(v, "createdAt__gt", &p.CreatedAtGt)
	setString(v, "createdAt__gte", &p.CreatedAtGte)
	setString(v, "createdAt__between", &p.CreatedAtBetween)
	return v
}

// -- API methods --

// ListUsers returns a paginated list of users visible to the authenticated
// user, filtered by the optional params.
//
// API: GET /web/api/v2.1/users
// Required permission: Users.view
//
// Pass nil for params to use the API defaults (limit 10, no filters).
func (c *Client) ListUsers(ctx context.Context, params *ListUsersParams) ([]User, *Pagination, error) {
	var p url.Values
	if params != nil {
		p = params.values()
	}
	var users []User
	pag, err := c.get(ctx, "/users", p, &users)
	if err != nil {
		return nil, nil, err
	}
	return users, pag, nil
}

// GetUser returns the user with the given userID.
//
// API: GET /web/api/v2.1/users/{user_id}
// Required permission: Users.view
func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	var user User
	_, err := c.get(ctx, fmt.Sprintf("/users/%s", userID), nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user account.
//
// API: POST /web/api/v2.1/users
// Required permission: Users.create
//
// Email, FullName, and Scope are required.  Scope must be one of "site",
// "account", or "tenant".  Assign roles via ScopeRoles (preferred) or the
// deprecated SiteRoles / TenantRoles fields.  If automatic onboarding is
// enabled for the tenant, omit Password — the user will receive an invitation
// email instead.
func (c *Client) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	var user User
	_, err := c.post(ctx, "/users", req, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates mutable fields on an existing user.
//
// API: PUT /web/api/v2.1/users/{user_id}
// Required permission: Users.update
//
// Scope is required in req.Data.  Only non-zero fields are applied.  To
// change a password, set both CurrentPassword and Password.  Role assignments
// are replaced in full when ScopeRoles is non-empty.
func (c *Client) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (*User, error) {
	var user User
	_, err := c.put(ctx, fmt.Sprintf("/users/%s", userID), req, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser permanently deletes the user with the given userID.
//
// API: DELETE /web/api/v2.1/users/{user_id}
// Required permission: Users.delete
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	_, err := c.delete(ctx, fmt.Sprintf("/users/%s", userID), nil, nil)
	return err
}

// BulkDeleteUsers permanently deletes all users that match the filter in req.
// Use with caution — this operation is irreversible.  Narrow the target set
// with req.Filter.IDs, req.Filter.Email, or scope-based criteria before
// calling.
//
// API: POST /web/api/v2.1/users/delete-users
// Required permission: Users.delete
func (c *Client) BulkDeleteUsers(ctx context.Context, req BulkUsersActionRequest) error {
	_, err := c.post(ctx, "/users/delete-users", req, nil)
	return err
}

// GetUserAPITokenDetails returns the creation and expiry timestamps of the API
// token belonging to the specified user.  It does not return the token value.
//
// API: GET /web/api/v2.1/users/{user_id}/api-token-details
// Required permission: Users.view
func (c *Client) GetUserAPITokenDetails(ctx context.Context, userID string) (*APITokenDetail, error) {
	var detail APITokenDetail
	_, err := c.get(ctx, fmt.Sprintf("/users/%s/api-token-details", userID), nil, &detail)
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
func (c *Client) GetAPITokenDetailsByToken(ctx context.Context, req GetAPITokenDetailsRequest) (*APITokenDetail, error) {
	var detail APITokenDetail
	_, err := c.post(ctx, "/users/api-token-details", req, &detail)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

// GenerateAPIToken generates a new API token for the currently authenticated
// user, replacing any existing token.  The plaintext token is returned in
// the response and cannot be retrieved again; store it securely.
//
// API: POST /web/api/v2.1/users/generate-api-token
// Required permission: (self — no special permission beyond being logged in)
func (c *Client) GenerateAPIToken(ctx context.Context, req GenerateAPITokenRequest) (*APITokenResponse, error) {
	var token APITokenResponse
	_, err := c.post(ctx, "/users/generate-api-token", req, &token)
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
func (c *Client) RevokeAPIToken(ctx context.Context, req UserIDRequest) error {
	_, err := c.post(ctx, "/users/revoke-api-token", req, nil)
	return err
}

// GenerateIFrameToken generates a short-lived token that allows embedding a
// SentinelOne console view in an iframe.  Scope the token to an account or
// site by setting req.Data.AccountID or req.Data.SiteID.
//
// API: POST /web/api/v2.1/users/generate-iframe-token
// Required permission: Users.generateIframeToken
func (c *Client) GenerateIFrameToken(ctx context.Context, req IFrameUserRequest) (*IFrameTokenResponse, error) {
	var token IFrameTokenResponse
	_, err := c.post(ctx, "/users/generate-iframe-token", req, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// Enable2FA enables the two-factor authentication requirement for the user
// identified by req.Data.UserID.  The user must subsequently complete
// enrollment via [Client.Enroll2FA].
//
// API: POST /web/api/v2.1/users/2fa/enable
// Required permission: Users.twoFa
func (c *Client) Enable2FA(ctx context.Context, req UserIDRequest) error {
	_, err := c.post(ctx, "/users/2fa/enable", req, nil)
	return err
}

// Disable2FA removes the two-factor authentication requirement for the user
// identified by req.Data.UserID.
//
// API: POST /web/api/v2.1/users/2fa/disable
// Required permission: Users.twoFa
func (c *Client) Disable2FA(ctx context.Context, req UserIDRequest) error {
	_, err := c.post(ctx, "/users/2fa/disable", req, nil)
	return err
}

// Enroll2FA initiates the 2FA enrollment flow for one or more users.  The
// response contains a TOTP secret and QR code URL that the user should scan
// with an authenticator app.
//
// API: POST /web/api/v2.1/users/enroll-2fa
// Required permission: Users.twoFa
func (c *Client) Enroll2FA(ctx context.Context, req UserIDsRequest) (*EnrollTFAResponse, error) {
	var resp EnrollTFAResponse
	_, err := c.post(ctx, "/users/enroll-2fa", req, &resp)
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
func (c *Client) Reset2FA(ctx context.Context, req ResetTFARequest) error {
	_, err := c.post(ctx, "/users/reset-2fa", req, nil)
	return err
}

// Delete2FA removes a user's 2FA configuration entirely, including any enrolled
// device.  After this call, 2FA is no longer configured for the user.
//
// API: POST /web/api/v2.1/users/delete-2fa
// Required permission: Users.twoFa
func (c *Client) Delete2FA(ctx context.Context, req DeleteTFARequest) error {
	_, err := c.post(ctx, "/users/delete-2fa", req, nil)
	return err
}

// ChangePassword changes the password for the currently authenticated user.
// CurrentPassword is required for verification; if 2FA is enabled, TwoFACode
// must also be provided.
//
// API: POST /web/api/v2.1/users/change-password
// Required permission: (self — no special permission beyond being logged in)
func (c *Client) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	_, err := c.post(ctx, "/users/change-password", req, nil)
	return err
}

// EnableApp enables or disables an integrated application for users.
// Set req.Data.Enable to true to enable the app, false to disable it.
//
// API: POST /web/api/v2.1/users/enable-app
// Required permission: Users.enableApp
func (c *Client) EnableApp(ctx context.Context, req EnableAppRequest) error {
	_, err := c.post(ctx, "/users/enable-app", req, nil)
	return err
}

// RequestApp submits a request for the authenticated user to access an
// application.  The response URL redirects to the app's access-request flow.
//
// API: POST /web/api/v2.1/users/request-app
// Required permission: (self — no special permission beyond being logged in)
func (c *Client) RequestApp(ctx context.Context, req RequestAppRequest) (*RequestAppResponse, error) {
	var resp RequestAppResponse
	_, err := c.post(ctx, "/users/request-app", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ElevateSession re-verifies the caller's identity by requiring their password,
// upgrading the current session token to an elevated-privilege token.  Use
// this before sensitive operations that require step-up authentication.
//
// API: POST /web/api/v2.1/users/auth/elevate
// Required permission: (self — no special permission beyond being logged in)
func (c *Client) ElevateSession(ctx context.Context, req ElevateSessionRequest) (*ElevateSessionResponse, error) {
	var resp ElevateSessionResponse
	_, err := c.post(ctx, "/users/auth/elevate", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// AuthEULA records that the authenticated user has accepted the End User
// License Agreement.  This is typically called once after first login when
// the server responds with a status indicating EULA acceptance is required.
//
// API: POST /web/api/v2.1/users/auth/eula
func (c *Client) AuthEULA(ctx context.Context) error {
	_, err := c.post(ctx, "/users/auth/eula", struct{}{}, nil)
	return err
}

// AuthApp completes an OAuth-style app authorization flow by exchanging an
// authorization code for a session token.
//
// API: POST /web/api/v2.1/users/auth/app
func (c *Client) AuthApp(ctx context.Context, req AuthAppRequest) (*LoginResponse, error) {
	var resp LoginResponse
	_, err := c.post(ctx, "/users/auth/app", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// IsTenantAdmin reports whether the currently authenticated user holds a
// tenant-admin role.  Returns true only when that role is confirmed.
//
// API: GET /web/api/v2.1/users/tenant-admin-auth-check
func (c *Client) IsTenantAdmin(ctx context.Context) (bool, error) {
	var result struct {
		IsAdmin bool `json:"isAdmin"`
	}
	_, err := c.get(ctx, "/users/tenant-admin-auth-check", nil, &result)
	if err != nil {
		return false, err
	}
	return result.IsAdmin, nil
}

// IsRSAuth reports whether the currently authenticated user has Remote Shell
// (RS) authorization.
//
// API: GET /web/api/v2.1/users/rs-auth-check
func (c *Client) IsRSAuth(ctx context.Context) (bool, error) {
	var result struct {
		IsAuth bool `json:"isAuth"`
	}
	_, err := c.get(ctx, "/users/rs-auth-check", nil, &result)
	if err != nil {
		return false, err
	}
	return result.IsAuth, nil
}

// IsViewerAuth reports whether the currently authenticated user holds at
// least a viewer-level role.
//
// API: GET /web/api/v2.1/users/viewer-auth-check
func (c *Client) IsViewerAuth(ctx context.Context) (bool, error) {
	var result struct {
		IsAuth bool `json:"isAuth"`
	}
	_, err := c.get(ctx, "/users/viewer-auth-check", nil, &result)
	if err != nil {
		return false, err
	}
	return result.IsAuth, nil
}

// Login authenticates with a username (email) and password, returning a session
// token on success.  If the account requires 2FA, the response Status will be
// "2fa_required" and TwoFAMethod will indicate which method to use; pass the
// resulting token and a TOTP code to [Client.LoginContinue].
//
// The returned token can be used to construct an authenticated [Client] with
// [NewClient].
//
// API: POST /web/api/v2.1/users/login
func (c *Client) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	var resp LoginResponse
	_, err := c.post(ctx, "/users/login", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// LoginContinue completes a multi-step login by providing the second factor.
// Populate req.Data.Token with the token from [Client.Login] and
// req.Data.Code with the user's TOTP or SMS code.
//
// API: POST /web/api/v2.1/users/login-continue
func (c *Client) LoginContinue(ctx context.Context, req LoginContinueRequest) (*LoginContinueResponse, error) {
	var resp LoginContinueResponse
	_, err := c.post(ctx, "/users/login-continue", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// LoginByToken authenticates using a short-lived login token (e.g. one
// received in a password-reset or invitation email).
//
// API: GET /web/api/v2.1/users/login/by-token
func (c *Client) LoginByToken(ctx context.Context, token string) (*LoginResponse, error) {
	p := url.Values{}
	p.Set("token", token)
	var resp LoginResponse
	_, err := c.get(ctx, "/users/login/by-token", p, &resp)
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
func (c *Client) LoginByAPIToken(ctx context.Context, req LoginByAPITokenRequest) (*LoginResponse, error) {
	var resp LoginResponse
	_, err := c.post(ctx, "/users/login/by-api-token", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetPassword sets a new password for a user using a password-reset token.
// Token comes from a password-reset email link; the response contains a
// session token if the operation succeeds.
//
// API: POST /web/api/v2.1/users/login/set-password
func (c *Client) SetPassword(ctx context.Context, req SetPasswordRequest) (*SetPasswordResponse, error) {
	var resp SetPasswordResponse
	_, err := c.post(ctx, "/users/login/set-password", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendResetPasswordEmail triggers a password-reset email for all users that
// match req.Filter.  Typically used with req.Filter.Email or req.Filter.IDs
// to target a single user.
//
// API: POST /web/api/v2.1/users/login/send-reset-password-email
// Required permission: Users.resetPassword
func (c *Client) SendResetPasswordEmail(ctx context.Context, req SendResetPasswordRequest) error {
	_, err := c.post(ctx, "/users/login/send-reset-password-email", req, nil)
	return err
}

// ForceResetPasswordOnLogin flags all users matching req.Filter so that they
// are required to choose a new password at their next login attempt.
//
// API: POST /web/api/v2.1/users/login/force-reset-password-on-login
// Required permission: Users.resetPassword
func (c *Client) ForceResetPasswordOnLogin(ctx context.Context, req ForceResetPasswordRequest) error {
	_, err := c.post(ctx, "/users/login/force-reset-password-on-login", req, nil)
	return err
}

// Logout invalidates the current session token.  After this call, the client's
// token can no longer be used to authenticate requests.
//
// API: POST /web/api/v2.1/users/logout
func (c *Client) Logout(ctx context.Context) error {
	_, err := c.post(ctx, "/users/logout", struct{}{}, nil)
	return err
}

// LoginSSO initiates a SAML2 SSO login flow and returns the redirect URL that
// the user's browser should be sent to in order to complete authentication
// with the identity provider.
//
// API: GET /web/api/v2.1/users/login/sso-saml2
func (c *Client) LoginSSO(ctx context.Context) (string, error) {
	var result struct {
		RedirectURL string `json:"redirectUrl,omitempty"`
	}
	_, err := c.get(ctx, "/users/login/sso-saml2", nil, &result)
	if err != nil {
		return "", err
	}
	return result.RedirectURL, nil
}

// LoginSSOForScope initiates a SAML2 SSO login scoped to a specific account or
// site (identified by scopeID), returning the identity-provider redirect URL.
//
// API: POST /web/api/v2.1/users/login/sso-saml2/{scope_id}
func (c *Client) LoginSSOForScope(ctx context.Context, scopeID string) (string, error) {
	var result struct {
		RedirectURL string `json:"redirectUrl,omitempty"`
	}
	_, err := c.post(ctx, fmt.Sprintf("/users/login/sso-saml2/%s", scopeID), struct{}{}, &result)
	if err != nil {
		return "", err
	}
	return result.RedirectURL, nil
}

// SSOReAuth initiates a SAML2 re-authentication challenge for the current
// session, typically triggered when an elevated-privilege action is required.
// Returns the identity-provider redirect URL.
//
// API: GET /web/api/v2.1/users/sso-saml2/re-auth
func (c *Client) SSOReAuth(ctx context.Context) (string, error) {
	var result struct {
		RedirectURL string `json:"redirectUrl,omitempty"`
	}
	_, err := c.get(ctx, "/users/sso-saml2/re-auth", nil, &result)
	if err != nil {
		return "", err
	}
	return result.RedirectURL, nil
}

// OnboardingVerify validates the onboarding token sent to a new user's email
// and optionally sets their initial password.  On success the response
// contains a session token.
//
// API: POST /web/api/v2.1/users/onboarding/verify
func (c *Client) OnboardingVerify(ctx context.Context, req OnboardingVerifyRequest) (*LoginResponse, error) {
	var resp LoginResponse
	_, err := c.post(ctx, "/users/onboarding/verify", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// OnboardingValidateToken checks whether the given onboarding token is still
// valid (not expired, not already used).  Returns true if valid.
//
// API: GET /web/api/v2.1/users/onboarding/validate-token
func (c *Client) OnboardingValidateToken(ctx context.Context, token string) (bool, error) {
	p := url.Values{}
	p.Set("token", token)
	var result struct {
		Valid bool `json:"valid,omitempty"`
	}
	_, err := c.get(ctx, "/users/onboarding/validate-token", p, &result)
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
func (c *Client) OnboardingSendVerificationEmail(ctx context.Context, req SendVerificationEmailRequest) error {
	_, err := c.post(ctx, "/users/onboarding/send-verification-email", req, nil)
	return err
}
