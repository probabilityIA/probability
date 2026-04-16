package app

import (
	"context"
	"fmt"
)

// RemovePermissionFromRole elimina un permiso específico de un rol
func (uc *RoleUseCase) RemovePermissionFromRole(ctx context.Context, roleID uint, permissionID uint) error {
	uc.log.Info().
		Uint("role_id", roleID).
		Uint("permission_id", permissionID).
		Msg("Eliminando permiso del rol")

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

	// Eliminar permiso usando el repositorio
	err = uc.repository.RemovePermissionFromRole(ctx, roleID, permissionID)
	if err != nil {
		uc.log.Error().
			Err(err).
			Uint("role_id", roleID).
			Uint("permission_id", permissionID).
			Msg("Error al eliminar permiso del rol")
		return err
	}

	uc.log.Info().
		Uint("role_id", roleID).
		Uint("permission_id", permissionID).
		Msg("Permiso eliminado exitosamente del rol")

	return nil
}
