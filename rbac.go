package sentinelone

import (
	"context"
	"fmt"
	"net/url"
)

// -- RBAC types --

// RolePermissionEntry represents a single permission within a role page.
type RolePermissionEntry struct {
	Identifier            string   `json:"identifier,omitempty"`
	Title                 string   `json:"title,omitempty"`
	Value                 bool     `json:"value,omitempty"`
	Description           string   `json:"description,omitempty"`
	AdditionalDescription string   `json:"additionalDescription,omitempty"`
	Type                  string   `json:"type,omitempty"`
	DependsOn             []string `json:"dependsOn,omitempty"`
	GroupName             string   `json:"groupName,omitempty"`
	DisabledReason        string   `json:"disabledReason,omitempty"`
	DisabledReasonCode    string   `json:"disabledReasonCode,omitempty"`
}

// RolePage represents a permissions page/section within a role definition.
type RolePage struct {
	Name        string                `json:"name,omitempty"`
	Identifier  string                `json:"identifier,omitempty"`
	Permissions []RolePermissionEntry `json:"permissions,omitempty"`
}

// Role represents a SentinelOne RBAC role.
type Role struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
	Creator        string `json:"creator,omitempty"`
	CreatorID      string `json:"creatorId,omitempty"`
	UpdatedBy      string `json:"updatedBy,omitempty"`
	UpdatedByID    string `json:"updatedById,omitempty"`
	UsersInRoles   int    `json:"usersInRoles,omitempty"`
	ScopeID        string `json:"scopeId,omitempty"`
	Scope          string `json:"scope,omitempty"`
	PredefinedRole bool   `json:"predefinedRole,omitempty"`
	AccountName    string `json:"accountName,omitempty"`
	SiteName       string `json:"siteName,omitempty"`
}

// RoleWithPermissions extends Role with its full permission page list.
type RoleWithPermissions struct {
	Role
	Pages []RolePage `json:"pages,omitempty"`
}

// -- Request types --

// RoleScopeFilter specifies the scope for RBAC operations.
type RoleScopeFilter struct {
	AccountIDs []string `json:"accountIds,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
	Tenant     *bool    `json:"tenant,omitempty"`
}

// CreateRoleRequest is the request body for POST /rbac/role.
type CreateRoleRequest struct {
	Data   CreateRoleData  `json:"data"`
	Filter RoleScopeFilter `json:"filter"`
}

// CreateRoleData holds the fields for creating a role.
type CreateRoleData struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permissionIds,omitempty"`
}

// UpdateRoleRequest is the request body for PUT /rbac/role/{id}.
type UpdateRoleRequest struct {
	Data   UpdateRoleData   `json:"data"`
	Filter *RoleScopeFilter `json:"filter,omitempty"`
}

// UpdateRoleData holds the fields for updating a role.
type UpdateRoleData struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permissionIds,omitempty"`
}

// DeleteRoleRequest is the request body for DELETE /rbac/role/{id}.
type DeleteRoleRequest struct {
	Data *struct{} `json:"data,omitempty"`
}

// -- Filter params --

// ListRolesParams contains query parameters for GET /web/api/v2.1/rbac/roles.
// All fields are optional; zero values are omitted from the request.
//
//   - AccountIDs / SiteIDs / GroupIDs: scope the listing to specific containers.
//   - Tenant: when true, list only tenant-scoped roles.
//   - PredefinedRole: true to list only system roles; false for custom roles only.
//   - IncludeParents: include roles from parent scopes (default true in the API).
//   - IncludeChildren: include roles from child scopes (default false in the API).
//   - CreatedAt* / UpdatedAt*: timestamp range filters using the __lt/__gt/__lte/
//     __gte/__between suffixes that the API accepts.
type ListRolesParams struct {
	ListParams
	AccountIDs       []string
	SiteIDs          []string
	GroupIDs         []string
	TenancyIDs       []string
	Tenant           *bool
	Name             string
	Query            string
	IDs              []string
	Creator          string
	CreatorID        string
	UpdatedBy        string
	UpdatedByID      string
	Description      string
	AccountName      string
	SiteName         string
	IncludeParents   *bool
	IncludeChildren  *bool
	PredefinedRole   *bool
	CreatedAt        string
	UpdatedAt        string
	CreatedAtLt      string
	CreatedAtGt      string
	CreatedAtLte     string
	CreatedAtGte     string
	CreatedAtBetween string
	UpdatedAtLt      string
	UpdatedAtGt      string
	UpdatedAtLte     string
	UpdatedAtGte     string
	UpdatedAtBetween string
}

func (p *ListRolesParams) values() url.Values {
	v := p.ListParams.values()
	setStringSlice(v, "accountIds", p.AccountIDs)
	setStringSlice(v, "siteIds", p.SiteIDs)
	setStringSlice(v, "groupIds", p.GroupIDs)
	setStringSlice(v, "tenancyIds", p.TenancyIDs)
	setBool(v, "tenant", p.Tenant)
	setString(v, "name", &p.Name)
	setString(v, "query", &p.Query)
	setStringSlice(v, "ids", p.IDs)
	setString(v, "creator", &p.Creator)
	setString(v, "creatorId", &p.CreatorID)
	setString(v, "updatedBy", &p.UpdatedBy)
	setString(v, "updatedById", &p.UpdatedByID)
	setString(v, "description", &p.Description)
	setString(v, "accountName", &p.AccountName)
	setString(v, "siteName", &p.SiteName)
	setBool(v, "includeParents", p.IncludeParents)
	setBool(v, "includeChildren", p.IncludeChildren)
	setBool(v, "predefinedRole", p.PredefinedRole)
	setString(v, "createdAt", &p.CreatedAt)
	setString(v, "updatedAt", &p.UpdatedAt)
	setString(v, "createdAt__lt", &p.CreatedAtLt)
	setString(v, "createdAt__gt", &p.CreatedAtGt)
	setString(v, "createdAt__lte", &p.CreatedAtLte)
	setString(v, "createdAt__gte", &p.CreatedAtGte)
	setString(v, "createdAt__between", &p.CreatedAtBetween)
	setString(v, "updatedAt__lt", &p.UpdatedAtLt)
	setString(v, "updatedAt__gt", &p.UpdatedAtGt)
	setString(v, "updatedAt__lte", &p.UpdatedAtLte)
	setString(v, "updatedAt__gte", &p.UpdatedAtGte)
	setString(v, "updatedAt__between", &p.UpdatedAtBetween)
	return v
}

// GetRoleTemplateParams scopes the template returned by [Client.GetRoleTemplate].
// At least one of AccountIDs, SiteIDs, GroupIDs, or Tenant should be set so
// the API returns the permission pages applicable to that scope level.
type GetRoleTemplateParams struct {
	AccountIDs []string
	SiteIDs    []string
	GroupIDs   []string
	Tenant     *bool
}

func (p *GetRoleTemplateParams) values() url.Values {
	v := url.Values{}
	setStringSlice(v, "accountIds", p.AccountIDs)
	setStringSlice(v, "siteIds", p.SiteIDs)
	setStringSlice(v, "groupIds", p.GroupIDs)
	setBool(v, "tenant", p.Tenant)
	return v
}

// GetRolePermissionsParams provides optional scope/filter context for
// [Client.GetRole].  Pass nil to fetch the role's permissions without
// additional scope filtering.
type GetRolePermissionsParams struct {
	AccountIDs       []string
	SiteIDs          []string
	GroupIDs         []string
	Tenant           *bool
	Name             string
	Query            string
	CreatedAtLt      string
	CreatedAtGt      string
	CreatedAtLte     string
	CreatedAtGte     string
	CreatedAtBetween string
	UpdatedAtLt      string
	UpdatedAtGt      string
	UpdatedAtLte     string
	UpdatedAtGte     string
	UpdatedAtBetween string
}

func (p *GetRolePermissionsParams) values() url.Values {
	v := url.Values{}
	setStringSlice(v, "accountIds", p.AccountIDs)
	setStringSlice(v, "siteIds", p.SiteIDs)
	setStringSlice(v, "groupIds", p.GroupIDs)
	setBool(v, "tenant", p.Tenant)
	setString(v, "name", &p.Name)
	setString(v, "query", &p.Query)
	setString(v, "createdAt__lt", &p.CreatedAtLt)
	setString(v, "createdAt__gt", &p.CreatedAtGt)
	setString(v, "createdAt__lte", &p.CreatedAtLte)
	setString(v, "createdAt__gte", &p.CreatedAtGte)
	setString(v, "createdAt__between", &p.CreatedAtBetween)
	setString(v, "updatedAt__lt", &p.UpdatedAtLt)
	setString(v, "updatedAt__gt", &p.UpdatedAtGt)
	setString(v, "updatedAt__lte", &p.UpdatedAtLte)
	setString(v, "updatedAt__gte", &p.UpdatedAtGte)
	setString(v, "updatedAt__between", &p.UpdatedAtBetween)
	return v
}

// -- API methods --

// ListRoles returns a paginated list of RBAC roles visible to the authenticated
// user, filtered by the optional params.
//
// API: GET /web/api/v2.1/rbac/roles
// Required permission: Roles.view
//
// By default the API includes parent-scope roles.  Set IncludeChildren true to
// also surface roles defined in child scopes (accounts, sites).
func (c *Client) ListRoles(ctx context.Context, params *ListRolesParams) ([]Role, *Pagination, error) {
	var p url.Values
	if params != nil {
		p = params.values()
	}
	var roles []Role
	pag, err := c.get(ctx, "/rbac/roles", p, &roles)
	if err != nil {
		return nil, nil, err
	}
	return roles, pag, nil
}

// GetRoleTemplate returns the blank permission-page structure used as a starting
// point for creating a new role in the given scope.  Each [RolePage] lists the
// available [RolePermissionEntry] items with their default values.  Pass these
// permission identifiers to [Client.CreateRole] via CreateRoleData.PermissionIDs.
//
// API: GET /web/api/v2.1/rbac/role
// Required permission: Roles.create (or Roles.view for read access to the template)
func (c *Client) GetRoleTemplate(ctx context.Context, params *GetRoleTemplateParams) (*RoleWithPermissions, error) {
	var p url.Values
	if params != nil {
		p = params.values()
	}
	var role RoleWithPermissions
	_, err := c.get(ctx, "/rbac/role", p, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRole returns the role identified by roleID, including its full list of
// permission pages and individual permission values.  Pass nil for params when
// no additional scope filtering is required.
//
// API: GET /web/api/v2.1/rbac/role/{role_id}
// Required permission: Roles.view
func (c *Client) GetRole(ctx context.Context, roleID string, params *GetRolePermissionsParams) (*RoleWithPermissions, error) {
	var p url.Values
	if params != nil {
		p = params.values()
	}
	var role RoleWithPermissions
	_, err := c.get(ctx, fmt.Sprintf("/rbac/role/%s", roleID), p, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateRole creates a new custom RBAC role in the scope specified by
// req.Filter.  Name and Description in req.Data are required.  Provide the
// permission identifiers from [Client.GetRoleTemplate] in PermissionIDs to
// grant specific access; omitting PermissionIDs creates a role with no
// permissions enabled.
//
// API: POST /web/api/v2.1/rbac/role
// Required permission: Roles.create
func (c *Client) CreateRole(ctx context.Context, req CreateRoleRequest) (*Role, error) {
	var role Role
	_, err := c.post(ctx, "/rbac/role", req, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// UpdateRole replaces the name, description, and permission set of an existing
// custom RBAC role.  Predefined (system) roles cannot be updated.
//
// API: PUT /web/api/v2.1/rbac/role/{role_id}
// Required permission: Roles.update
func (c *Client) UpdateRole(ctx context.Context, roleID string, req UpdateRoleRequest) (*Role, error) {
	var role Role
	_, err := c.put(ctx, fmt.Sprintf("/rbac/role/%s", roleID), req, &role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// DeleteRole permanently removes a custom RBAC role.  Users assigned to this
// role will lose the associated permissions.  Predefined (system) roles cannot
// be deleted.
//
// API: DELETE /web/api/v2.1/rbac/role/{role_id}
// Required permission: Roles.delete
func (c *Client) DeleteRole(ctx context.Context, roleID string) error {
	_, err := c.delete(ctx, fmt.Sprintf("/rbac/role/%s", roleID), DeleteRoleRequest{}, nil)
	return err
}
