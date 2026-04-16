package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
)

// GetPermissionsByScopeID obtiene permisos por scope ID
func (uc *PermissionUseCase) GetPermissionsByScopeID(ctx context.Context, scopeID uint) ([]domain.PermissionDTO, error) {
	uc.logger.Info().Uint("scope_id", scopeID).Msg("Obteniendo permisos por scope ID")

	permissions, err := uc.repository.GetPermissionsByScopeID(ctx, scopeID)
	if err != nil {
		uc.logger.Error().Uint("scope_id", scopeID).Err(err).Msg("Error al obtener permisos por scope ID desde el repositorio")
		return nil, err
	}

	// Convertir entidades a DTOs
	permissionDTOs := make([]domain.PermissionDTO, len(permissions))
	for i, permission := range permissions {
		permissionDTOs[i] = entityToPermissionDTO(permission)
	}

	uc.logger.Info().Uint("scope_id", scopeID).Int("count", len(permissionDTOs)).Msg("Permisos por scope ID obtenidos exitosamente")
	return permissionDTOs, nil
}
