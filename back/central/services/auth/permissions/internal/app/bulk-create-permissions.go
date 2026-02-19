package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
)

// BulkCreatePermissions crea múltiples permisos en una sola operación
func (uc *PermissionUseCase) BulkCreatePermissions(ctx context.Context, dtos []domain.CreatePermissionDTO) ([]domain.BulkCreateResult, error) {
	results := make([]domain.BulkCreateResult, len(dtos))

	for i, dto := range dtos {
		result, err := uc.CreatePermission(ctx, dto)
		if err != nil {
			results[i] = domain.BulkCreateResult{Name: dto.Name, Success: false, Error: err.Error()}
		} else {
			results[i] = domain.BulkCreateResult{Name: dto.Name, Success: true, Message: result}
		}
	}

	return results, nil
}
