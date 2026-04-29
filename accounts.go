package sentinelone

import (
	"context"
	"net/url"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

const (
	accountsRootPath = "/accounts"
	accountsBasePath = accountsRootPath + "/"
)

// AccountsClient provides access to the Accounts API group.
// Access it via [Client.Accounts].
type AccountsClient struct{ c *Client }

// List returns a paginated list of accounts visible to the authenticated user,
// filtered by the optional params.
//
// API: GET /web/api/v2.1/accounts
// Required permission: Accounts.view
//
// Pass nil for params to use the API defaults (limit 10, no filters).
// Use [types.Pagination].NextCursor for subsequent pages.
func (a *AccountsClient) List(
	ctx context.Context,
	params *ListAccountsParams,
) ([]types.Account, *types.Pagination, error) {
	var paramVals url.Values
	if params != nil {
		paramVals = params.values()
	}

	var accounts []types.Account

	pag, err := a.c.get(ctx, accountsRootPath, paramVals, &accounts)
	if err != nil {
		return nil, nil, err
	}

	return accounts, pag, nil
}

// Get returns the account with the given accountID.
//
// API: GET /web/api/v2.1/accounts/{account_id}
// Required permission: Accounts.view
func (a *AccountsClient) Get(ctx context.Context, accountID string) (*types.Account, error) {
	var account types.Account

	_, err := a.c.get(ctx, accountsBasePath+accountID, nil, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// Create creates a new account under the tenant.
//
// API: POST /web/api/v2.1/accounts
// Required permission: Accounts.create
//
// The Name field in req.Data is required.  AccountType ("Trial" or "Paid")
// and Expiration are optional; omit Expiration or set UnlimitedExpiration true
// for accounts that should never expire.
func (a *AccountsClient) Create(ctx context.Context, req CreateAccountRequest) (*types.Account, error) {
	var account types.Account

	_, err := a.c.post(ctx, accountsRootPath, req, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// Update updates mutable fields on an existing account.
//
// API: PUT /web/api/v2.1/accounts/{account_id}
// Required permission: Accounts.update
//
// Only fields set in req.Data are sent; zero values are omitted.  To update
// the security policy in the same call, populate req.Data.Policy.
func (a *AccountsClient) Update(
	ctx context.Context,
	accountID string,
	req UpdateAccountRequest,
) (*types.Account, error) {
	var account types.Account

	_, err := a.c.put(ctx, accountsBasePath+accountID, req, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetPolicy returns the active security policy for an account.
// If the account inherits from the global tenant policy, the inherited values
// are reflected in the returned [types.Policy].
//
// API: GET /web/api/v2.1/accounts/{account_id}/policy
// Required permission: Accounts.editPolicy
func (a *AccountsClient) GetPolicy(ctx context.Context, accountID string) (*types.Policy, error) {
	var policy types.Policy

	_, err := a.c.get(ctx, accountsBasePath+accountID+"/policy", nil, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// UpdatePolicy replaces the security policy for an account.
// Setting Inherits to true on the account itself (via [AccountsClient.Update]) causes
// it to follow the global policy; this call overrides individual fields.
//
// API: PUT /web/api/v2.1/accounts/{account_id}/policy
// Required permission: Accounts.editPolicy
func (a *AccountsClient) UpdatePolicy(
	ctx context.Context,
	accountID string,
	req UpdatePolicyRequest,
) (*types.Policy, error) {
	var policy types.Policy

	_, err := a.c.put(ctx, accountsBasePath+accountID+"/policy", req, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// RevertPolicy removes any account-level policy overrides, causing the account
// to inherit the global tenant policy again.
//
// API: PUT /web/api/v2.1/accounts/{account_id}/revert-policy
// Required permission: Accounts.revertPolicy
func (a *AccountsClient) RevertPolicy(ctx context.Context, accountID string) error {
	_, err := a.c.put(ctx, accountsBasePath+accountID+"/revert-policy", RevertPolicyRequest{}, nil)

	return err
}

// Reactivate transitions an expired account back to the "active" state.
// Provide a new Expiration timestamp or set UnlimitedExpiration to true in
// req.Data; at least one must be supplied.
//
// API: PUT /web/api/v2.1/accounts/{account_id}/reactivate
// Required permission: Accounts.reactivate
func (a *AccountsClient) Reactivate(ctx context.Context, accountID string, req ReactivateAccountRequest) (
	*types.Account, error,
) {
	var account types.Account

	_, err := a.c.put(ctx, accountsBasePath+accountID+"/reactivate", req, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// ExpireNow immediately transitions an active account to the "expired" state
// without waiting for its scheduled expiration date.
//
// API: POST /web/api/v2.1/accounts/{account_id}/expire-now
// Required permission: Accounts.expire
func (a *AccountsClient) ExpireNow(ctx context.Context, accountID string) error {
	_, err := a.c.post(ctx, accountsBasePath+accountID+"/expire-now", struct{}{}, nil)

	return err
}

// GetUninstallPasswordMetadata returns metadata about the current uninstall
// password for an account — creation time, creator, and expiry — without
// revealing the password itself.  Use [AccountsClient.ViewUninstallPassword]
// when the actual value is needed.
//
// API: GET /web/api/v2.1/accounts/{account_id}/uninstall-password/metadata
// Required permission: Accounts.uninstallPassword.view
func (a *AccountsClient) GetUninstallPasswordMetadata(ctx context.Context, accountID string) (
	*types.UninstallPasswordMetadata, error,
) {
	var meta types.UninstallPasswordMetadata

	_, err := a.c.get(ctx, accountsBasePath+accountID+"/uninstall-password/metadata", nil, &meta)
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
func (a *AccountsClient) ViewUninstallPassword(
	ctx context.Context,
	accountID string,
) (*types.UninstallPassword, error) {
	var pass types.UninstallPassword

	_, err := a.c.get(ctx, accountsBasePath+accountID+"/uninstall-password/view", nil, &pass)
	if err != nil {
		return nil, err
	}

	return &pass, nil
}

// GenerateUninstallPassword creates a new uninstall password for an account,
// replacing any previously generated password.  The new plaintext password is
// returned in the response; store it securely as it cannot be retrieved again.
//
// API: POST /web/api/v2.1/accounts/{account_id}/uninstall-password/generate
// Required permission: Accounts.uninstallPassword.manage
func (a *AccountsClient) GenerateUninstallPassword(
	ctx context.Context,
	accountID string,
) (*types.UninstallPassword, error) {
	var pass types.UninstallPassword

	_, err := a.c.post(ctx, accountsBasePath+accountID+"/uninstall-password/generate",
		GenerateUninstallPasswordRequest{}, &pass)
	if err != nil {
		return nil, err
	}

	return &pass, nil
}

// RevokeUninstallPassword removes the current uninstall password from an
// account.  Agents protected by this password will no longer require it after
// the revocation propagates.
//
// API: POST /web/api/v2.1/accounts/{account_id}/uninstall-password/revoke
// Required permission: Accounts.uninstallPassword.manage
func (a *AccountsClient) RevokeUninstallPassword(ctx context.Context, accountID string) error {
	_, err := a.c.post(ctx, accountsBasePath+accountID+"/uninstall-password/revoke",
		RevokeUninstallPasswordRequest{}, nil)

	return err
}
