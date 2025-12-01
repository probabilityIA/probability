package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
)

// GetRolesByLevel obtiene roles por nivel
func (uc *RoleUseCase) GetRolesByLevel(ctx context.Context, filters domain.RoleFilters) ([]domain.RoleDTO, error) {
	level := 0
	if filters.Level != nil {
		level = *filters.Level
	}

	uc.log.Info().Int("level", level).Msg("Iniciando caso de uso: obtener roles por nivel")

	roles, err := uc.repository.GetRolesByLevel(ctx, level)
	if err != nil {
		uc.log.Error().Int("level", level).Err(err).Msg("Error al obtener roles por nivel desde el repositorio")
		return nil, err
	}

	// Convertir entidades a DTOs
	roleDTOs := make([]domain.RoleDTO, len(roles))
	for i, role := range roles {
		roleDTOs[i] = entityToRoleDTO(role)
	}

	uc.log.Info().Int("level", level).Int("count", len(roleDTOs)).Msg("Roles por nivel obtenidos exitosamente")
	return roleDTOs, nil
}
