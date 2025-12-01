package app

import (
	"context"
	"fmt"
)

// AssignPermissionsToRole asigna permisos a un rol
func (uc *RoleUseCase) AssignPermissionsToRole(ctx context.Context, roleID uint, permissionIDs []uint) error {
	uc.log.Info().
		Uint("role_id", roleID).
		Int("permission_count", len(permissionIDs)).
		Msg("Asignando permisos a rol")

	// Verificar que el rol existe
	role, err := uc.repository.GetRoleByID(ctx, roleID)
	if err != nil {
		uc.log.Error().
			Err(err).
			Uint("role_id", roleID).
			Msg("Error al verificar existencia del rol")
		return fmt.Errorf("rol no encontrado")
	}

	if role == nil {
		uc.log.Error().
			Uint("role_id", roleID).
			Msg("Rol no encontrado")
		return fmt.Errorf("rol no encontrado")
	}

	// Asignar permisos usando el repositorio
	err = uc.repository.AssignPermissionsToRole(ctx, roleID, permissionIDs)
	if err != nil {
		uc.log.Error().
			Err(err).
			Uint("role_id", roleID).
			Msg("Error al asignar permisos al rol")
		return err
	}

	uc.log.Info().
		Uint("role_id", roleID).
		Int("permission_count", len(permissionIDs)).
		Msg("Permisos asignados exitosamente al rol")

	return nil
}
