package domain

import "context"

type IPermissionRepository interface {
	GetPermissions(ctx context.Context, businessTypeID *uint, name *string, scopeID *uint) ([]Permission, error)
	GetPermissionByID(ctx context.Context, id uint) (*Permission, error)
	GetPermissionsByScopeID(ctx context.Context, scopeID uint) ([]Permission, error)
	GetPermissionsByResource(ctx context.Context, resource string) ([]Permission, error)
	PermissionExistsByName(ctx context.Context, name string) (bool, error)
	CreatePermission(ctx context.Context, permission Permission) (string, error)
	UpdatePermission(ctx context.Context, id uint, permission Permission) (string, error)
	DeletePermission(ctx context.Context, id uint) (string, error)
}
