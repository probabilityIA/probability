package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
)

// GetRolePermissions obtiene los permisos de un rol
func (uc *RoleUseCase) GetRolePermissions(ctx context.Context, roleID uint) ([]domain.Permission, error) {
	uc.log.Info().
		Uint("role_id", roleID).
		Msg("Obteniendo permisos del rol")

	// Verificar que el rol existe
	role, err := uc.repository.GetRoleByID(ctx, roleID)
	if err != nil {
		uc.log.Error().
			Err(err).
			Uint("role_id", roleID).
			Msg("Error al verificar existencia del rol")
		return nil, fmt.Errorf("rol no encontrado")
	}

	if role == nil {
		uc.log.Error().
			Uint("role_id", roleID).
			Msg("Rol no encontrado")
		return nil, fmt.Errorf("rol no encontrado")
	}

	// Obtener permisos usando el repositorio
	permissions, err := uc.repository.GetRolePermissions(ctx, roleID)
	if err != nil {
		uc.log.Error().
			Err(err).
			Uint("role_id", roleID).
			Msg("Error al obtener permisos del rol")
		return nil, err
	}

	uc.log.Info().
		Uint("role_id", roleID).
		Int("permission_count", len(permissions)).
		Msg("Permisos del rol obtenidos exitosamente")

	return permissions, nil
}
