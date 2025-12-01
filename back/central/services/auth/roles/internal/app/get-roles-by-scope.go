package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
)

// GetRolesByScopeID obtiene roles por scope ID
func (uc *RoleUseCase) GetRolesByScopeID(ctx context.Context, scopeID uint) ([]domain.RoleDTO, error) {
	uc.log.Info().Uint("scope_id", scopeID).Msg("Iniciando caso de uso: obtener roles por scope ID")

	roles, err := uc.repository.GetRolesByScopeID(ctx, scopeID)
	if err != nil {
		uc.log.Error().Uint("scope_id", scopeID).Err(err).Msg("Error al obtener roles por scope ID desde el repositorio")
		return nil, err
	}

	// Convertir entidades a DTOs
	roleDTOs := make([]domain.RoleDTO, len(roles))
	for i, role := range roles {
		roleDTOs[i] = entityToRoleDTO(role)
	}

	uc.log.Info().Uint("scope_id", scopeID).Int("count", len(roleDTOs)).Msg("Roles por scope obtenidos exitosamente")
	return roleDTOs, nil
}
