package domain

import "context"

type IRoleRepository interface {
	GetRoleByID(ctx context.Context, id uint) (*Role, error)
	AssignPermissionsToRole(ctx context.Context, roleID uint, permissionIDs []uint) error
	CreateRole(ctx context.Context, role CreateRoleDTO) (*Role, error)
	UpdateRole(ctx context.Context, id uint, role UpdateRoleDTO) (*Role, error)
	RoleExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]Permission, error)
	GetRolesByLevel(ctx context.Context, level int) ([]Role, error)
	GetRolesByScopeID(ctx context.Context, scopeID uint) ([]Role, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetSystemRoles(ctx context.Context) ([]Role, error)
	RemovePermissionFromRole(ctx context.Context, roleID uint, permissionID uint) error
}
