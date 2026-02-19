package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/domain"
)

// GetPermissions obtiene todos los permisos con filtros opcionales
func (uc *PermissionUseCase) GetPermissions(ctx context.Context, businessTypeID *uint, name *string, scopeID *uint, resource *string) ([]domain.PermissionDTO, error) {
	uc.logger.Info().Msg("Obteniendo todos los permisos")

	permissions, err := uc.repository.GetPermissions(ctx, businessTypeID, name, scopeID, resource)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error al obtener permisos desde el repositorio")
		return nil, err
	}

	// Convertir entidades a DTOs
	permissionDTOs := make([]domain.PermissionDTO, len(permissions))
	for i, permission := range permissions {
		permissionDTOs[i] = entityToPermissionDTO(permission)
	}

	uc.logger.Info().Int("count", len(permissionDTOs)).Msg("Permisos obtenidos exitosamente")
	return permissionDTOs, nil
}

// entityToPermissionDTO convierte una entidad Permission a PermissionDTO
func entityToPermissionDTO(permission domain.Permission) domain.PermissionDTO {
	return domain.PermissionDTO{
		ID:               permission.ID,
		Name:             permission.Name,
		Code:             permission.Code,
		Description:      permission.Description,
		Resource:         permission.Resource,
		Action:           permission.Action,
		ResourceID:       permission.ResourceID,
		ActionID:         permission.ActionID,
		ScopeID:          permission.ScopeID,
		ScopeName:        permission.ScopeName,
		ScopeCode:        permission.ScopeCode,
		BusinessTypeID:   permission.BusinessTypeID,
		BusinessTypeName: permission.BusinessTypeName,
	}
}
