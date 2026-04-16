package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
)

// AssignRoleToUserBusiness asigna o actualiza roles de un usuario en múltiples businesses
func (uc *UserUseCase) AssignRoleToUserBusiness(ctx context.Context, userID uint, assignments []domain.BusinessRoleAssignment) error {
	uc.log.Info().
		Uint("user_id", userID).
		Int("assignments_count", len(assignments)).
		Msg("Iniciando asignación de roles a usuario en businesses")

	if len(assignments) == 0 {
		uc.log.Error().Uint("user_id", userID).Msg("No se proporcionaron asignaciones")
		return fmt.Errorf("no se proporcionaron asignaciones")
	}

	// Ejecutar asignación en el repositorio
	// El repositorio valida:
	// - Que el usuario existe
	// - Que todos los businesses existen
	// - Que todos los roles existen
	// - Que cada rol es del mismo tipo de business que su business asociado
	// - Que el usuario está asociado a cada business en business_staff
	if err := uc.repository.AssignRoleToUserBusiness(ctx, userID, assignments); err != nil {
		uc.log.Error().
			Err(err).
			Uint("user_id", userID).
			Int("assignments_count", len(assignments)).
			Msg("Error al asignar roles a usuario en businesses")
		return fmt.Errorf("error al asignar roles: %w", err)
	}

	uc.log.Info().
		Uint("user_id", userID).
		Int("assignments_count", len(assignments)).
		Msg("Roles asignados exitosamente a usuario en businesses")

	return nil
}
