package sentinelone

import "github.com/s1buildpartners/sentinelone-go-sdk/types"

// StateActive, StateExpired, and StateDeleted are the valid values for State
// filter fields on accounts, sites, and related endpoints.
const (
	StateActive  = "active"
	StateExpired = "expired"
	StateDeleted = "deleted"
)

// UpdatePolicyRequest is the request body for PUT /accounts/{id}/policy or PUT /sites/{id}/policy.
type UpdatePolicyRequest struct {
	Data types.Policy `json:"data"`
}

// RevertPolicyRequest is the request body for PUT /accounts/{id}/revert-policy or PUT /sites/{id}/revert-policy.
type RevertPolicyRequest struct {
	Data *struct{} `json:"data,omitempty"`
}
