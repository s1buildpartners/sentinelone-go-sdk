package sentinelone

import (
	"context"
	"fmt"
	"net/url"
)

// -- Account types --

// AccountSKU represents a deprecated SKU entry on an account.
type AccountSKU struct {
	Type          string `json:"type,omitempty"`
	TotalLicenses int    `json:"totalLicenses,omitempty"`
	Unlimited     bool   `json:"unlimited,omitempty"`
	AgentsInSKU   int    `json:"agentsInSku,omitempty"`
}

// LicenseSurface is a surface within a license bundle.
type LicenseSurface struct {
	Count int    `json:"count,omitempty"`
	Name  string `json:"name,omitempty"`
}

// LicenseBundle is a product bundle entry in a license set.
type LicenseBundle struct {
	DisplayName   string           `json:"displayName,omitempty"`
	MajorVersion  int              `json:"majorVersion,omitempty"`
	MinorVersion  int              `json:"minorVersion,omitempty"`
	Name          string           `json:"name,omitempty"`
	Surfaces      []LicenseSurface `json:"surfaces,omitempty"`
	TotalSurfaces int              `json:"totalSurfaces,omitempty"`
}

// LicenseModule is an add-on module in a license set.
type LicenseModule struct {
	DisplayName  string `json:"displayName,omitempty"`
	MajorVersion int    `json:"majorVersion,omitempty"`
	Name         string `json:"name,omitempty"`
}

// LicenseSetting is a deprecated license setting entry.
type LicenseSetting struct {
	DisplayName              string `json:"displayName,omitempty"`
	SettingGroup             string `json:"settingGroup,omitempty"`
	SettingGroupDisplayName  string `json:"settingGroupDisplayName,omitempty"`
}

// Licenses holds the complete license information for an account or site.
type Licenses struct {
	Bundles  []LicenseBundle  `json:"bundles,omitempty"`
	Modules  []LicenseModule  `json:"modules,omitempty"`
	Settings []LicenseSetting `json:"settings,omitempty"`
}

// IRFields holds IR (Incident Response) contact and classification data.
type IRFields struct {
	CompanyName                  string  `json:"companyName"`
	ContactFirstName             string  `json:"contactFirstName"`
	ContactLastName              string  `json:"contactLastName"`
	ContactEmail                 string  `json:"contactEmail"`
	Region                       string  `json:"region,omitempty"`
	Country                      string  `json:"country"`
	City                         *string `json:"city,omitempty"`
	Postal                       *string `json:"postal,omitempty"`
	NumberOfEmployeesEndpoints   int     `json:"numberOfEmployeesEndpoints"`
	Industry                     string  `json:"industry"`
}

// PolicyEngines holds engine on/off states for a policy.
type PolicyEngines struct {
	Reputation             string `json:"reputation,omitempty"`
	PreExecution           string `json:"preExecution,omitempty"`
	PreExecutionSuspicious string `json:"preExecutionSuspicious,omitempty"`
	Executables            string `json:"executables,omitempty"`
	DataFiles              string `json:"dataFiles,omitempty"`
	Exploits               string `json:"exploits,omitempty"`
	Penetration            string `json:"penetration,omitempty"`
	PUP                    string `json:"pup,omitempty"`
}

// Policy represents a site or account security policy.
type Policy struct {
	NetworkQuarantineOn    *bool          `json:"networkQuarantineOn,omitempty"`
	AutoImmuneOn           *bool          `json:"autoImmuneOn,omitempty"`
	AutoDecommissionOn     *bool          `json:"autoDecommissionOn,omitempty"`
	IsDefault              *bool          `json:"isDefault,omitempty"`
	ResearchOn             *bool          `json:"researchOn,omitempty"`
	AutoMitigationAction   string         `json:"autoMitigationAction,omitempty"`
	AutoDecommissionDays   *int           `json:"autoDecommissionDays,omitempty"`
	MitigationMode         string         `json:"mitigationMode,omitempty"`
	CreatedAt              *string        `json:"createdAt,omitempty"`
	AgentNotification      *bool          `json:"agentNotification,omitempty"`
	Engines                *PolicyEngines `json:"engines,omitempty"`
	ScanNewAgents          *bool          `json:"scanNewAgents,omitempty"`
	MonitoringOn           *bool          `json:"monitoringOn,omitempty"`
	InheritedFrom          string         `json:"inheritedFrom,omitempty"`
	UpdatedAt              *string        `json:"updatedAt,omitempty"`
}

// Account represents a SentinelOne management account.
type Account struct {
	ID                  string      `json:"id,omitempty"`
	Name                string      `json:"name,omitempty"`
	IsDefault           bool        `json:"isDefault,omitempty"`
	AccountType         string      `json:"accountType,omitempty"`
	Expiration          *string     `json:"expiration,omitempty"`
	SalesforceID        string      `json:"salesforceId,omitempty"`
	ExternalID          string      `json:"externalId,omitempty"`
	State               string      `json:"state,omitempty"`
	CreatedAt           *string     `json:"createdAt,omitempty"`
	UpdatedAt           *string     `json:"updatedAt,omitempty"`
	UnlimitedExpiration bool        `json:"unlimitedExpiration,omitempty"`
	ActiveAgents        int         `json:"activeAgents,omitempty"`
	TotalLicenses       int         `json:"totalLicenses,omitempty"`
	NumberOfSites       int         `json:"numberOfSites,omitempty"`
	UsageType           string      `json:"usageType,omitempty"`
	BillingMode         string      `json:"billingMode,omitempty"`
	MakeSOCDefaultUI    bool        `json:"makeSocDefaultUi,omitempty"`
	Creator             string      `json:"creator,omitempty"`
	CreatorID           string      `json:"creatorId,omitempty"`
	Licenses            *Licenses   `json:"licenses,omitempty"`
	IRFields            *IRFields   `json:"irFields,omitempty"`
	SKUs                []AccountSKU `json:"skus,omitempty"`

	// Deprecated fields kept for backwards compatibility
	CoreSites           int  `json:"coreSites,omitempty"`
	ControlSites        int  `json:"controlSites,omitempty"`
	CompleteSites       int  `json:"completeSites,omitempty"`
	TotalCore           int  `json:"totalCore,omitempty"`
	TotalControl        int  `json:"totalControl,omitempty"`
	TotalComplete       int  `json:"totalComplete,omitempty"`
	UnlimitedCore       bool `json:"unlimitedCore,omitempty"`
	UnlimitedControl    bool `json:"unlimitedControl,omitempty"`
	UnlimitedComplete   bool `json:"unlimitedComplete,omitempty"`
	AgentsInCoreSKU     int  `json:"agentsInCoreSku,omitempty"`
	AgentsInControlSKU  int  `json:"agentsInControlSku,omitempty"`
	AgentsInCompleteSKU int  `json:"agentsInCompleteSku,omitempty"`
}

// UninstallPasswordMetadata holds metadata about an account's uninstall password.
type UninstallPasswordMetadata struct {
	CreatedAt      *string `json:"createdAt,omitempty"`
	CreatedBy      string  `json:"createdBy,omitempty"`
	CreatedByID    string  `json:"createdById,omitempty"`
	ExpiresAt      *string `json:"expiresAt,omitempty"`
	HasPassword    bool    `json:"hasPassword,omitempty"`
}

// UninstallPassword holds a generated uninstall password.
type UninstallPassword struct {
	Password string `json:"password,omitempty"`
}

// -- Request types --

// CreateAccountRequest is the request body for POST /accounts.
type CreateAccountRequest struct {
	Data CreateAccountData `json:"data"`
}

// CreateAccountData holds the fields for creating an account.
type CreateAccountData struct {
	Name        string  `json:"name"`
	AccountType string  `json:"accountType,omitempty"`
	Expiration  *string `json:"expiration,omitempty"`
	ExternalID  *string `json:"externalId,omitempty"`
	Inherits    *bool   `json:"inherits,omitempty"`
	Policy      *Policy `json:"policy,omitempty"`
}

// UpdateAccountRequest is the request body for PUT /accounts/{id}.
type UpdateAccountRequest struct {
	Data UpdateAccountData `json:"data"`
}

// UpdateAccountData holds the fields for updating an account.
type UpdateAccountData struct {
	Name                string   `json:"name,omitempty"`
	AccountType         string   `json:"accountType,omitempty"`
	Expiration          *string  `json:"expiration,omitempty"`
	ExternalID          *string  `json:"externalId,omitempty"`
	UnlimitedExpiration *bool    `json:"unlimitedExpiration,omitempty"`
	Inherits            *bool    `json:"inherits,omitempty"`
	Policy              *Policy  `json:"policy,omitempty"`
}

// ReactivateAccountRequest is the request body for PUT /accounts/{id}/reactivate.
type ReactivateAccountRequest struct {
	Data ReactivateAccountData `json:"data"`
}

// ReactivateAccountData holds reactivation parameters.
type ReactivateAccountData struct {
	Expiration          *string `json:"expiration,omitempty"`
	UnlimitedExpiration *bool   `json:"unlimitedExpiration,omitempty"`
}

// UpdatePolicyRequest is the request body for PUT /accounts/{id}/policy or PUT /sites/{id}/policy.
type UpdatePolicyRequest struct {
	Data Policy `json:"data"`
}

// RevertPolicyRequest is the request body for PUT /accounts/{id}/revert-policy.
type RevertPolicyRequest struct {
	Data *struct{} `json:"data,omitempty"`
}

// GenerateUninstallPasswordRequest is the request body for POST /accounts/{id}/uninstall-password/generate.
type GenerateUninstallPasswordRequest struct {
	Data *struct{} `json:"data,omitempty"`
}

// RevokeUninstallPasswordRequest is the request body for POST /accounts/{id}/uninstall-password/revoke.
type RevokeUninstallPasswordRequest struct {
	Data *struct{} `json:"data,omitempty"`
}

// -- Filter params for list operations --

// ListAccountsParams contains query parameters for GET /web/api/v2.1/accounts.
// All fields are optional; zero values are omitted from the request.
//
//   - IDs / AccountIDs: filter to specific account IDs (max 5,000 per request).
//   - Query: full-text search on account name.
//   - Name: exact name match.
//   - AccountType: "Trial" or "Paid".
//   - State: "active", "expired", or "deleted".
//   - States / StatesNin: include or exclude multiple states.
//   - Features: filter accounts that support a feature ("firewall-control", etc.).
//   - UsageType: "customer", "mssp", or "ir".
//   - BillingMode: "subscription" or "consumption".
//   - SKU / Module: filter by product SKU or module identifier.
//   - NameContains: partial-match filters on account name (multi-value OR).
type ListAccountsParams struct {
	ListParams
	IDs          []string
	AccountIDs   []string
	Query        string
	Name         string
	IsDefault    *bool
	AccountType  string // Trial, Paid
	State        string // active, expired, deleted
	States       []string
	StatesNin    []string
	Features     []string
	UsageType    string // customer, mssp, ir
	BillingMode  string // subscription, consumption
	SKU          string
	Module       string
	Expiration   string
	CreatedAt    string
	UpdatedAt    string
	NameContains []string
}

func (p *ListAccountsParams) values() url.Values {
	v := p.ListParams.values()
	setStringSlice(v, "ids", p.IDs)
	setStringSlice(v, "accountIds", p.AccountIDs)
	setString(v, "query", &p.Query)
	setString(v, "name", &p.Name)
	setBool(v, "isDefault", p.IsDefault)
	setString(v, "accountType", &p.AccountType)
	setString(v, "state", &p.State)
	setStringSlice(v, "states", p.States)
	setStringSlice(v, "statesNin", p.StatesNin)
	setStringSlice(v, "features", p.Features)
	setString(v, "usageType", &p.UsageType)
	setString(v, "billingMode", &p.BillingMode)
	setString(v, "sku", &p.SKU)
	setString(v, "module", &p.Module)
	setString(v, "expiration", &p.Expiration)
	setString(v, "createdAt", &p.CreatedAt)
	setString(v, "updatedAt", &p.UpdatedAt)
	setStringSlice(v, "name__contains", p.NameContains)
	return v
}

// -- API methods --

// ListAccounts returns a paginated list of accounts visible to the authenticated
// user, filtered by the optional params.
//
// API: GET /web/api/v2.1/accounts
// Required permission: Accounts.view
//
// Pass nil for params to use the API defaults (limit 10, no filters).
// Use [Pagination].NextCursor for subsequent pages.
func (c *Client) ListAccounts(ctx context.Context, params *ListAccountsParams) ([]Account, *Pagination, error) {
	var p url.Values
	if params != nil {
		p = params.values()
	}
	var accounts []Account
	pag, err := c.get(ctx, "/accounts", p, &accounts)
	if err != nil {
		return nil, nil, err
	}
	return accounts, pag, nil
}

// GetAccount returns the account with the given accountID.
//
// API: GET /web/api/v2.1/accounts/{account_id}
// Required permission: Accounts.view
func (c *Client) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	var account Account
	_, err := c.get(ctx, fmt.Sprintf("/accounts/%s", accountID), nil, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// CreateAccount creates a new account under the tenant.
//
// API: POST /web/api/v2.1/accounts
// Required permission: Accounts.create
//
// The Name field in req.Data is required.  AccountType ("Trial" or "Paid")
// and Expiration are optional; omit Expiration or set UnlimitedExpiration true
// for accounts that should never expire.
func (c *Client) CreateAccount(ctx context.Context, req CreateAccountRequest) (*Account, error) {
	var account Account
	_, err := c.post(ctx, "/accounts", req, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// UpdateAccount updates mutable fields on an existing account.
//
// API: PUT /web/api/v2.1/accounts/{account_id}
// Required permission: Accounts.update
//
// Only fields set in req.Data are sent; zero values are omitted.  To update
// the security policy in the same call, populate req.Data.Policy.
func (c *Client) UpdateAccount(ctx context.Context, accountID string, req UpdateAccountRequest) (*Account, error) {
	var account Account
	_, err := c.put(ctx, fmt.Sprintf("/accounts/%s", accountID), req, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountPolicy returns the active security policy for an account.
// If the account inherits from the global tenant policy, the inherited values
// are reflected in the returned [Policy].
//
// API: GET /web/api/v2.1/accounts/{account_id}/policy
// Required permission: Accounts.editPolicy
func (c *Client) GetAccountPolicy(ctx context.Context, accountID string) (*Policy, error) {
	var policy Policy
	_, err := c.get(ctx, fmt.Sprintf("/accounts/%s/policy", accountID), nil, &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// UpdateAccountPolicy replaces the security policy for an account.
// Setting Inherits to true on the account itself (via [UpdateAccount]) causes
// it to follow the global policy; this call overrides individual fields.
//
// API: PUT /web/api/v2.1/accounts/{account_id}/policy
// Required permission: Accounts.editPolicy
func (c *Client) UpdateAccountPolicy(ctx context.Context, accountID string, req UpdatePolicyRequest) (*Policy, error) {
	var policy Policy
	_, err := c.put(ctx, fmt.Sprintf("/accounts/%s/policy", accountID), req, &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// RevertAccountPolicy removes any account-level policy overrides, causing the
// account to inherit the global tenant policy again.
//
// API: PUT /web/api/v2.1/accounts/{account_id}/revert-policy
// Required permission: Accounts.revertPolicy
func (c *Client) RevertAccountPolicy(ctx context.Context, accountID string) error {
	_, err := c.put(ctx, fmt.Sprintf("/accounts/%s/revert-policy", accountID), RevertPolicyRequest{}, nil)
	return err
}

// ReactivateAccount transitions an expired account back to the "active" state.
// Provide a new Expiration timestamp or set UnlimitedExpiration to true in
// req.Data; at least one must be supplied.
//
// API: PUT /web/api/v2.1/accounts/{account_id}/reactivate
// Required permission: Accounts.reactivate
func (c *Client) ReactivateAccount(ctx context.Context, accountID string, req ReactivateAccountRequest) (*Account, error) {
	var account Account
	_, err := c.put(ctx, fmt.Sprintf("/accounts/%s/reactivate", accountID), req, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// ExpireAccountNow immediately transitions an active account to the "expired"
// state without waiting for its scheduled expiration date.
//
// API: POST /web/api/v2.1/accounts/{account_id}/expire-now
// Required permission: Accounts.expire
func (c *Client) ExpireAccountNow(ctx context.Context, accountID string) error {
	_, err := c.post(ctx, fmt.Sprintf("/accounts/%s/expire-now", accountID), struct{}{}, nil)
	return err
}

// GetUninstallPasswordMetadata returns metadata about the current uninstall
// password for an account — creation time, creator, and expiry — without
// revealing the password itself.  Use [Client.ViewUninstallPassword] when the
// actual value is needed.
//
// API: GET /web/api/v2.1/accounts/{account_id}/uninstall-password/metadata
// Required permission: Accounts.uninstallPassword.view
func (c *Client) GetUninstallPasswordMetadata(ctx context.Context, accountID string) (*UninstallPasswordMetadata, error) {
	var meta UninstallPasswordMetadata
	_, err := c.get(ctx, fmt.Sprintf("/accounts/%s/uninstall-password/metadata", accountID), nil, &meta)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

// ViewUninstallPassword returns the plaintext uninstall password currently
// configured for an account.  This is a sensitive operation; the caller must
// hold the appropriate permission.
//
// API: GET /web/api/v2.1/accounts/{account_id}/uninstall-password/view
// Required permission: Accounts.uninstallPassword.view
func (c *Client) ViewUninstallPassword(ctx context.Context, accountID string) (*UninstallPassword, error) {
	var pw UninstallPassword
	_, err := c.get(ctx, fmt.Sprintf("/accounts/%s/uninstall-password/view", accountID), nil, &pw)
	if err != nil {
		return nil, err
	}
	return &pw, nil
}

// GenerateUninstallPassword creates a new uninstall password for an account,
// replacing any previously generated password.  The new plaintext password is
// returned in the response; store it securely as it cannot be retrieved again.
//
// API: POST /web/api/v2.1/accounts/{account_id}/uninstall-password/generate
// Required permission: Accounts.uninstallPassword.manage
func (c *Client) GenerateUninstallPassword(ctx context.Context, accountID string) (*UninstallPassword, error) {
	var pw UninstallPassword
	_, err := c.post(ctx, fmt.Sprintf("/accounts/%s/uninstall-password/generate", accountID), GenerateUninstallPasswordRequest{}, &pw)
	if err != nil {
		return nil, err
	}
	return &pw, nil
}

// RevokeUninstallPassword removes the current uninstall password from an
// account.  Agents protected by this password will no longer require it after
// the revocation propagates.
//
// API: POST /web/api/v2.1/accounts/{account_id}/uninstall-password/revoke
// Required permission: Accounts.uninstallPassword.manage
func (c *Client) RevokeUninstallPassword(ctx context.Context, accountID string) error {
	_, err := c.post(ctx, fmt.Sprintf("/accounts/%s/uninstall-password/revoke", accountID), RevokeUninstallPasswordRequest{}, nil)
	return err
}
