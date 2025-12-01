package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
)

// GetPermissionsByResource obtiene permisos por recurso
func (uc *PermissionUseCase) GetPermissionsByResource(ctx context.Context, resource string) ([]domain.PermissionDTO, error) {
	uc.logger.Info().Str("resource", resource).Msg("Obteniendo permisos por recurso")

	permissions, err := uc.repository.GetPermissionsByResource(ctx, resource)
	if err != nil {
		uc.logger.Error().Str("resource", resource).Err(err).Msg("Error al obtener permisos por recurso desde el repositorio")
		return nil, err
	}

	// Convertir entidades a DTOs
	permissionDTOs := make([]domain.PermissionDTO, len(permissions))
	for i, permission := range permissions {
		permissionDTOs[i] = entityToPermissionDTO(permission)
	}

	uc.logger.Info().Str("resource", resource).Int("count", len(permissionDTOs)).Msg("Permisos por recurso obtenidos exitosamente")
	return permissionDTOs, nil
}
