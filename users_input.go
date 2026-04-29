package sentinelone

import (
	"net/url"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// CreateUserRequest is the request body for POST /users.
type CreateUserRequest struct {
	Data CreateUserData `json:"data"`
}

// CreateUserData holds the fields for creating a user.
type CreateUserData struct {
	Email               string                `json:"email"`
	FullName            string                `json:"fullName"`
	Scope               string                `json:"scope"`
	Password            string                `json:"password,omitempty"`
	TwoFAEnabled        *bool                 `json:"twoFaEnabled,omitempty"`
	ForceCredentialAuth *bool                 `json:"forceCredentialAuth,omitempty"`
	ScopeRoles          []types.UserScopeRole `json:"scopeRoles,omitempty"`
	SiteRoles           []types.UserSiteRole  `json:"siteRoles,omitempty"`
	TenantRoles         []string              `json:"tenantRoles,omitempty"`
}

// UpdateUserRequest is the request body for PUT /users/{id}.
type UpdateUserRequest struct {
	Data UpdateUserData `json:"data"`
}

// UpdateUserData holds the fields for updating a user.
type UpdateUserData struct {
	Scope               string                `json:"scope"`
	ID                  string                `json:"id,omitempty"`
	Email               string                `json:"email,omitempty"`
	FullName            string                `json:"fullName,omitempty"`
	Password            string                `json:"password,omitempty"`
	CurrentPassword     string                `json:"currentPassword,omitempty"`
	TwoFAEnabled        *bool                 `json:"twoFaEnabled,omitempty"`
	TwoFACode           string                `json:"twoFaCode,omitempty"`
	CanGenerateAPIToken *bool                 `json:"canGenerateApiToken,omitempty"`
	ScopeRoles          []types.UserScopeRole `json:"scopeRoles,omitempty"`
	SiteRoles           []types.UserSiteRole  `json:"siteRoles,omitempty"`
	TenantRoles         []string              `json:"tenantRoles,omitempty"`
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
	IDs                      []string `json:"ids,omitempty"`
	Email                    string   `json:"email,omitempty"`
	EmailContains            []string `json:"email__contains,omitempty"`
	FullName                 string   `json:"fullName,omitempty"`
	FullNameContains         []string `json:"fullName__contains,omitempty"`
	Source                   string   `json:"source,omitempty"`
	Sources                  []string `json:"sources,omitempty"`
	TwoFAEnabled             *bool    `json:"twoFaEnabled,omitempty"`
	TwoFAStatus              string   `json:"twoFaStatus,omitempty"`
	TwoFAStatuses            []string `json:"twoFaStatuses,omitempty"`
	EmailVerified            *bool    `json:"emailVerified,omitempty"`
	EmailReadOnly            *bool    `json:"emailReadOnly,omitempty"`
	FullNameReadOnly         *bool    `json:"fullNameReadOnly,omitempty"`
	GroupsReadOnly           *bool    `json:"groupsReadOnly,omitempty"`
	RoleIDs                  []string `json:"roleIds,omitempty"`
	CanGenerateAPIToken      *bool    `json:"canGenerateApiToken,omitempty"`
	HasValidAPIToken         *bool    `json:"hasValidApiToken,omitempty"`
	Query                    string   `json:"query,omitempty"`
	LastActivationLt         string   `json:"lastActivation__lt,omitempty"`
	LastActivationLte        string   `json:"lastActivation__lte,omitempty"`
	LastActivationGt         string   `json:"lastActivation__gt,omitempty"`
	LastActivationGte        string   `json:"lastActivation__gte,omitempty"`
	LastActivationBetween    string   `json:"lastActivation__between,omitempty"`
	APITokenExpiresAtLt      string   `json:"apiTokenExpiresAt__lt,omitempty"`
	APITokenExpiresAtLte     string   `json:"apiTokenExpiresAt__lte,omitempty"`
	APITokenExpiresAtGt      string   `json:"apiTokenExpiresAt__gt,omitempty"`
	APITokenExpiresAtGte     string   `json:"apiTokenExpiresAt__gte,omitempty"`
	APITokenExpiresAtBetween string   `json:"apiTokenExpiresAt__between,omitempty"`
	CreatedAtLt              string   `json:"createdAt__lt,omitempty"`
	CreatedAtLte             string   `json:"createdAt__lte,omitempty"`
	CreatedAtGt              string   `json:"createdAt__gt,omitempty"`
	CreatedAtGte             string   `json:"createdAt__gte,omitempty"`
	CreatedAtBetween         string   `json:"createdAt__between,omitempty"`
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

func (p *ListUsersParams) values() url.Values { //nolint:funlen
	vals := p.ListParams.values()
	setStringSlice(vals, "siteIds", p.SiteIDs)
	setStringSlice(vals, "accountIds", p.AccountIDs)
	setStringSlice(vals, "ids", p.IDs)
	setStringSlice(vals, "roleIds", p.RoleIDs)
	setString(vals, "source", &p.Source)
	setStringSlice(vals, "sources", p.Sources)
	setString(vals, "email", &p.Email)
	setStringSlice(vals, "email__contains", p.EmailContains)
	setBool(vals, "emailReadOnly", p.EmailReadOnly)
	setString(vals, "fullName", &p.FullName)
	setStringSlice(vals, "fullName__contains", p.FullNameContains)
	setBool(vals, "fullNameReadOnly", p.FullNameReadOnly)
	setBool(vals, "twoFaEnabled", p.TwoFAEnabled)
	setString(vals, "twoFaStatus", &p.TwoFAStatus)
	setStringSlice(vals, "twoFaStatuses", p.TwoFAStatuses)
	setString(vals, "primaryTwoFaMethod", &p.PrimaryTwoFAMethod)
	setBool(vals, "emailVerified", p.EmailVerified)
	setBool(vals, "canGenerateApiToken", p.CanGenerateAPIToken)
	setBool(vals, "hasValidApiToken", p.HasValidAPIToken)
	setString(vals, "query", &p.Query)
	setBool(vals, "groupsReadOnly", p.GroupsReadOnly)
	setString(vals, "firstLogin", &p.FirstLogin)
	setString(vals, "lastLogin", &p.LastLogin)
	setString(vals, "dateJoined", &p.DateJoined)
	setString(vals, "lastActivation__lt", &p.LastActivationLt)
	setString(vals, "lastActivation__lte", &p.LastActivationLte)
	setString(vals, "lastActivation__gt", &p.LastActivationGt)
	setString(vals, "lastActivation__gte", &p.LastActivationGte)
	setString(vals, "lastActivation__between", &p.LastActivationBetween)
	setString(vals, "apiTokenExpiresAt__lt", &p.APITokenExpiresAtLt)
	setString(vals, "apiTokenExpiresAt__lte", &p.APITokenExpiresAtLte)
	setString(vals, "apiTokenExpiresAt__gt", &p.APITokenExpiresAtGt)
	setString(vals, "apiTokenExpiresAt__gte", &p.APITokenExpiresAtGte)
	setString(vals, "apiTokenExpiresAt__between", &p.APITokenExpiresAtBetween)
	setString(vals, "createdAt__lt", &p.CreatedAtLt)
	setString(vals, "createdAt__lte", &p.CreatedAtLte)
	setString(vals, "createdAt__gt", &p.CreatedAtGt)
	setString(vals, "createdAt__gte", &p.CreatedAtGte)
	setString(vals, "createdAt__between", &p.CreatedAtBetween)

	return vals
}
