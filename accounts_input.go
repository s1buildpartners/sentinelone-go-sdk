package sentinelone

import (
	"net/url"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// CreateAccountRequest is the request body for POST /accounts.
type CreateAccountRequest struct {
	Data CreateAccountData `json:"data"`
}

// CreateAccountData holds the fields for creating an account.
type CreateAccountData struct {
	Name        string        `json:"name"`
	AccountType string        `json:"accountType,omitempty"`
	Expiration  *string       `json:"expiration,omitempty"`
	ExternalID  *string       `json:"externalId,omitempty"`
	Inherits    *bool         `json:"inherits,omitempty"`
	Policy      *types.Policy `json:"policy,omitempty"`
}

// UpdateAccountRequest is the request body for PUT /accounts/{id}.
type UpdateAccountRequest struct {
	Data UpdateAccountData `json:"data"`
}

// UpdateAccountData holds the fields for updating an account.
type UpdateAccountData struct {
	Name                string        `json:"name,omitempty"`
	AccountType         string        `json:"accountType,omitempty"`
	Expiration          *string       `json:"expiration,omitempty"`
	ExternalID          *string       `json:"externalId,omitempty"`
	UnlimitedExpiration *bool         `json:"unlimitedExpiration,omitempty"`
	Inherits            *bool         `json:"inherits,omitempty"`
	Policy              *types.Policy `json:"policy,omitempty"`
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

// GenerateUninstallPasswordRequest is the request body for POST /accounts/{id}/uninstall-password/generate.
type GenerateUninstallPasswordRequest struct {
	Data *struct{} `json:"data,omitempty"`
}

// RevokeUninstallPasswordRequest is the request body for POST /accounts/{id}/uninstall-password/revoke.
type RevokeUninstallPasswordRequest struct {
	Data *struct{} `json:"data,omitempty"`
}

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
	vals := p.ListParams.values()
	setStringSlice(vals, "ids", p.IDs)
	setStringSlice(vals, "accountIds", p.AccountIDs)
	setString(vals, "query", &p.Query)
	setString(vals, "name", &p.Name)
	setBool(vals, "isDefault", p.IsDefault)
	setString(vals, "accountType", &p.AccountType)
	setString(vals, "state", &p.State)
	setStringSlice(vals, "states", p.States)
	setStringSlice(vals, "statesNin", p.StatesNin)
	setStringSlice(vals, "features", p.Features)
	setString(vals, "usageType", &p.UsageType)
	setString(vals, "billingMode", &p.BillingMode)
	setString(vals, "sku", &p.SKU)
	setString(vals, "module", &p.Module)
	setString(vals, "expiration", &p.Expiration)
	setString(vals, "createdAt", &p.CreatedAt)
	setString(vals, "updatedAt", &p.UpdatedAt)
	setStringSlice(vals, "name__contains", p.NameContains)

	return vals
}
