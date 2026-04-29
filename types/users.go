package types

// UserAPIToken holds summary API token metadata embedded in a [User] response.
// For the standalone token-details endpoint, see [APITokenDetail].
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
	ID                   string          `json:"id,omitempty"`
	Source               string          `json:"source,omitempty"`
	GlobalUserID         *string         `json:"globalUserId,omitempty"`
	GlobalOrganizationID *string         `json:"globalOrganizationId,omitempty"`
	Email                *string         `json:"email,omitempty"`
	EmailReadOnly        bool            `json:"emailReadOnly,omitempty"`
	FullName             *string         `json:"fullName,omitempty"`
	FullNameReadOnly     bool            `json:"fullNameReadOnly,omitempty"`
	FirstLogin           string          `json:"firstLogin,omitempty"`
	LastLogin            string          `json:"lastLogin,omitempty"`
	DateJoined           string          `json:"dateJoined,omitempty"`
	TwoFAEnabled         bool            `json:"twoFaEnabled,omitempty"`
	TwoFAEnabledReadOnly bool            `json:"twoFaEnabledReadOnly,omitempty"`
	TwoFAConfigured      bool            `json:"twoFaConfigured,omitempty"`
	TwoFAStatus          string          `json:"twoFaStatus,omitempty"`
	PrimaryTwoFAMethod   string          `json:"primaryTwoFaMethod,omitempty"`
	EmailVerified        bool            `json:"emailVerified,omitempty"`
	Scope                string          `json:"scope"`
	APIToken             *UserAPIToken   `json:"apiToken,omitempty"`
	AgreedEULA           bool            `json:"agreedEula,omitempty"`
	AgreementURL         string          `json:"agreementUrl,omitempty"`
	CanGenerateAPIToken  bool            `json:"canGenerateApiToken,omitempty"`
	IsSystem             bool            `json:"isSystem,omitempty"`
	ScopeRoles           []UserScopeRole `json:"scopeRoles,omitempty"`
	SiteRoles            []UserSiteRole  `json:"siteRoles,omitempty"`
	TenantRoles          []string        `json:"tenantRoles,omitempty"`

	// Deprecated
	LowestRole       string `json:"lowestRole,omitempty"`
	GroupsReadOnly   bool   `json:"groupsReadOnly,omitempty"`
	AllowRemoteShell bool   `json:"allowRemoteShell,omitempty"`
}

// LoginResponse holds the result of a login call.
// When 2FA is required, Status is "2fa_required" and Token is a temporary
// continuation token to pass to the login-continue endpoint; TwoFAMethod
// indicates the required method ("totp", "sms", etc.).  On full success,
// Status is "active" and Token is the authenticated session token.
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

// APITokenDetail holds the creation and expiry timestamps of a user's API token,
// as returned by the token-details endpoints.  It does not include the token
// value itself.
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

// ElevateSessionResponse holds the upgraded session token returned after a
// successful session-elevation (step-up authentication) call.
type ElevateSessionResponse struct {
	Token string `json:"token,omitempty"`
}

// RequestAppResponse holds the app-access request result.
type RequestAppResponse struct {
	URL string `json:"url,omitempty"`
}

// SetPasswordResponse holds the authenticated session token returned after a
// password is successfully set via a password-reset token.
type SetPasswordResponse struct {
	Token string `json:"token,omitempty"`
}
