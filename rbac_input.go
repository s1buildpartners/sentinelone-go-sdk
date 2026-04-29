package sentinelone

import "net/url"

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
	Data struct{} `json:"data"`
}

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
	vals := p.ListParams.values()
	setStringSlice(vals, "accountIds", p.AccountIDs)
	setStringSlice(vals, "siteIds", p.SiteIDs)
	setStringSlice(vals, "groupIds", p.GroupIDs)
	setStringSlice(vals, "tenancyIds", p.TenancyIDs)
	setBool(vals, "tenant", p.Tenant)
	setString(vals, "name", &p.Name)
	setString(vals, "query", &p.Query)
	setStringSlice(vals, "ids", p.IDs)
	setString(vals, "creator", &p.Creator)
	setString(vals, "creatorId", &p.CreatorID)
	setString(vals, "updatedBy", &p.UpdatedBy)
	setString(vals, "updatedById", &p.UpdatedByID)
	setString(vals, "description", &p.Description)
	setString(vals, "accountName", &p.AccountName)
	setString(vals, "siteName", &p.SiteName)
	setBool(vals, "includeParents", p.IncludeParents)
	setBool(vals, "includeChildren", p.IncludeChildren)
	setBool(vals, "predefinedRole", p.PredefinedRole)
	setString(vals, "createdAt", &p.CreatedAt)
	setString(vals, "updatedAt", &p.UpdatedAt)
	setString(vals, "createdAt__lt", &p.CreatedAtLt)
	setString(vals, "createdAt__gt", &p.CreatedAtGt)
	setString(vals, "createdAt__lte", &p.CreatedAtLte)
	setString(vals, "createdAt__gte", &p.CreatedAtGte)
	setString(vals, "createdAt__between", &p.CreatedAtBetween)
	setString(vals, "updatedAt__lt", &p.UpdatedAtLt)
	setString(vals, "updatedAt__gt", &p.UpdatedAtGt)
	setString(vals, "updatedAt__lte", &p.UpdatedAtLte)
	setString(vals, "updatedAt__gte", &p.UpdatedAtGte)
	setString(vals, "updatedAt__between", &p.UpdatedAtBetween)

	return vals
}

// GetRoleTemplateParams scopes the template returned by [RBACClient.GetTemplate].
// At least one of AccountIDs, SiteIDs, GroupIDs, or Tenant should be set so
// the API returns the permission pages applicable to that scope level.
type GetRoleTemplateParams struct {
	AccountIDs []string
	SiteIDs    []string
	GroupIDs   []string
	Tenant     *bool
}

func (p *GetRoleTemplateParams) values() url.Values {
	vals := url.Values{}
	setStringSlice(vals, "accountIds", p.AccountIDs)
	setStringSlice(vals, "siteIds", p.SiteIDs)
	setStringSlice(vals, "groupIds", p.GroupIDs)
	setBool(vals, "tenant", p.Tenant)

	return vals
}

// GetRolePermissionsParams provides optional scope/filter context for
// [RBACClient.Get].  Pass nil to fetch the role's permissions without
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
	vals := url.Values{}
	setStringSlice(vals, "accountIds", p.AccountIDs)
	setStringSlice(vals, "siteIds", p.SiteIDs)
	setStringSlice(vals, "groupIds", p.GroupIDs)
	setBool(vals, "tenant", p.Tenant)
	setString(vals, "name", &p.Name)
	setString(vals, "query", &p.Query)
	setString(vals, "createdAt__lt", &p.CreatedAtLt)
	setString(vals, "createdAt__gt", &p.CreatedAtGt)
	setString(vals, "createdAt__lte", &p.CreatedAtLte)
	setString(vals, "createdAt__gte", &p.CreatedAtGte)
	setString(vals, "createdAt__between", &p.CreatedAtBetween)
	setString(vals, "updatedAt__lt", &p.UpdatedAtLt)
	setString(vals, "updatedAt__gt", &p.UpdatedAtGt)
	setString(vals, "updatedAt__lte", &p.UpdatedAtLte)
	setString(vals, "updatedAt__gte", &p.UpdatedAtGte)
	setString(vals, "updatedAt__between", &p.UpdatedAtBetween)

	return vals
}
