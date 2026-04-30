package sentinelone

import (
	"context"
	"net/url"

	"github.com/s1buildpartners/sentinelone-go-sdk/types"
)

// RBACClient provides access to the RBAC (roles) API group.
// Access it via [Client.RBAC].
type RBACClient struct{ c *Client }

// List returns a paginated list of RBAC roles visible to the authenticated
// user, filtered by the optional params.
//
// API: GET /web/api/v2.1/rbac/roles
//
// By default the API includes parent-scope roles.  Set IncludeChildren true to
// also surface roles defined in child scopes (accounts, sites).
func (r *RBACClient) List(ctx context.Context, params *ListRolesParams) ([]types.Role, *types.Pagination, error) {
	var paramVals url.Values
	if params != nil {
		paramVals = params.values()
	}

	var roles []types.Role

	pag, err := r.c.get(ctx, "/rbac/roles", paramVals, &roles)
	if err != nil {
		return nil, nil, err
	}

	return roles, pag, nil
}

// GetTemplate returns the blank permission-page structure used as a starting
// point for creating a new role in the given scope.  Each [types.RolePage] lists the
// available [types.RolePermissionEntry] items with their default values.  Pass these
// permission identifiers to [RBACClient.Create] via CreateRoleData.PermissionIDs.
//
// API: GET /web/api/v2.1/rbac/role
func (r *RBACClient) GetTemplate(
	ctx context.Context,
	params *GetRoleTemplateParams,
) (*types.RoleWithPermissions, error) {
	var paramVals url.Values
	if params != nil {
		paramVals = params.values()
	}

	var role types.RoleWithPermissions

	_, err := r.c.get(ctx, "/rbac/role", paramVals, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// Get returns the role identified by roleID, including its full list of
// permission pages and individual permission values.  Pass nil for params when
// no additional scope filtering is required.
//
// API: GET /web/api/v2.1/rbac/role/{role_id}
func (r *RBACClient) Get(
	ctx context.Context,
	roleID string,
	params *GetRolePermissionsParams,
) (*types.RoleWithPermissions, error) {
	var paramVals url.Values
	if params != nil {
		paramVals = params.values()
	}

	var role types.RoleWithPermissions

	_, err := r.c.get(ctx, "/rbac/role/"+roleID, paramVals, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// Create creates a new custom RBAC role in the scope specified by
// req.Filter.  Name and Description in req.Data are required.  Provide the
// permission identifiers from [RBACClient.GetTemplate] in PermissionIDs to
// grant specific access; omitting PermissionIDs creates a role with no
// permissions enabled.
//
// API: POST /web/api/v2.1/rbac/role
func (r *RBACClient) Create(ctx context.Context, req CreateRoleRequest) (*types.Role, error) {
	var role types.Role

	_, err := r.c.post(ctx, "/rbac/role", req, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// Update replaces the name, description, and permission set of an existing
// custom RBAC role.  Predefined (system) roles cannot be updated.
//
// API: PUT /web/api/v2.1/rbac/role/{role_id}
func (r *RBACClient) Update(ctx context.Context, roleID string, req UpdateRoleRequest) (*types.Role, error) {
	var role types.Role

	_, err := r.c.put(ctx, "/rbac/role/"+roleID, req, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// Delete permanently removes a custom RBAC role.  Users assigned to this role
// will lose the associated permissions.  Predefined (system) roles cannot be
// deleted.
//
// API: DELETE /web/api/v2.1/rbac/role/{role_id}
func (r *RBACClient) Delete(ctx context.Context, roleID string) error {
	_, err := r.c.delete(ctx, "/rbac/role/"+roleID, DeleteRoleRequest{})

	return err
}
