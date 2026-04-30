package sentinelone

import (
	"context"
	"net/url"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

const (
	agentsRootPath = "/agents"
)

// AgentsClient provides access to the Agents API group.
// Access it via [Client.Agents].
type AgentsClient struct{ c *Client }

// List returns a paginated list of agents visible to the authenticated user,
// filtered by the optional params.
//
// API: GET /web/api/v2.1/agents
// Required permission: Agents.view
//
// Pass nil for params to use the API defaults (limit 10, no filters).
// Use [types.Pagination].NextCursor for subsequent pages.
func (a *AgentsClient) List(
	ctx context.Context,
	params *ListAgentsParams,
) ([]types.Agent, *types.Pagination, error) {
	var paramVals url.Values

	if params != nil {
		paramVals = params.values()
	}

	var agents []types.Agent

	pag, err := a.c.get(ctx, agentsRootPath, paramVals, &agents)
	if err != nil {
		return nil, nil, err
	}

	return agents, pag, nil
}

// Count returns the total number of agents matching the given filter.
//
// API: GET /web/api/v2.1/agents/count
// Required permission: Agents.view
//
// Pass nil to count all agents visible to the caller.
func (a *AgentsClient) Count(ctx context.Context, params *ListAgentsParams) (int, error) {
	var paramVals url.Values

	if params != nil {
		paramVals = params.values()
	}

	var result struct {
		Total int `json:"total"`
	}

	_, err := a.c.get(ctx, agentsRootPath+"/count", paramVals, &result)
	if err != nil {
		return 0, err
	}

	return result.Total, nil
}
