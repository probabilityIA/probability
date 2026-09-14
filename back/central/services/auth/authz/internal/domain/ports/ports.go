package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/entities"
)

type IRepository interface {
	GetStaffRole(ctx context.Context, userID uint, businessID *uint) (*entities.StaffRole, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]entities.RolePermission, error)
	GetBusinessResourceConfig(ctx context.Context, businessID uint) ([]entities.ResourceConfig, error)
	GetBusinessName(ctx context.Context, businessID uint) (string, error)
}

type IModuleAccess interface {
	GetAccessibleModules(ctx context.Context, businessID uint) ([]string, error)
}

type IAccessCache interface {
	Get(ctx context.Context, userID, businessID uint) (*entities.Access, bool)
	Set(ctx context.Context, access *entities.Access)
}
