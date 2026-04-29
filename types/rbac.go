package types

// RolePermissionEntry represents a single permission toggle within a role page.
// Identifier is the opaque string to pass to CreateRoleData.PermissionIDs or
// UpdateRoleData.PermissionIDs to grant that permission.  Value indicates
// whether the permission is currently enabled for the role.  DependsOn lists
// Identifiers of other permissions that must also be enabled for this one to
// take effect.
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
