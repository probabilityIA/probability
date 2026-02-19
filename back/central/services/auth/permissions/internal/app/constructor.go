package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// PermissionUseCase implementa los casos de uso para permisos
type Iapp interface {
	GetPermissions(ctx context.Context, businessTypeID *uint, name *string, scopeID *uint, resource *string) ([]domain.PermissionDTO, error)
	GetPermissionByID(ctx context.Context, id uint) (*domain.PermissionDTO, error)
	CreatePermission(ctx context.Context, permission domain.CreatePermissionDTO) (string, error)
	BulkCreatePermissions(ctx context.Context, permissions []domain.CreatePermissionDTO) ([]domain.BulkCreateResult, error)
	UpdatePermission(ctx context.Context, id uint, permission domain.UpdatePermissionDTO) (string, error)
	DeletePermission(ctx context.Context, id uint) (string, error)
	GetPermissionsByResource(ctx context.Context, resource string) ([]domain.PermissionDTO, error)
	GetPermissionsByScopeID(ctx context.Context, scopeID uint) ([]domain.PermissionDTO, error)
}

type PermissionUseCase struct {
	repository domain.IPermissionRepository
	logger     log.ILogger
}

// NewPermissionUseCase crea una nueva instancia del caso de uso de permisos
func New(repository domain.IPermissionRepository, logger log.ILogger) Iapp {
	return &PermissionUseCase{
		repository: repository,
		logger:     logger,
	}
}
