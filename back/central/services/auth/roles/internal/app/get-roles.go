package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
)

// GetRoles obtiene todos los roles
func (uc *RoleUseCase) GetRoles(ctx context.Context, filters domain.RoleFilters) ([]domain.RoleDTO, error) {
	uc.log.Info().Msg("Iniciando caso de uso: obtener todos los roles")

	roles, err := uc.repository.GetRoles(ctx)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error al obtener roles desde el repositorio")
		return nil, err
	}

	// Aplicar filtros
	filteredRoles := []domain.Role{}
	for _, role := range roles {
		// Filtrar por BusinessTypeID
		if filters.BusinessTypeID != nil && role.BusinessTypeID != *filters.BusinessTypeID {
			continue
		}

		// Filtrar por ScopeID
		if filters.ScopeID != nil && role.ScopeID != *filters.ScopeID {
			continue
		}

		// Filtrar por IsSystem
		if filters.IsSystem != nil && role.IsSystem != *filters.IsSystem {
			continue
		}

		// Filtrar por Name (búsqueda parcial case-insensitive)
		if filters.Name != nil && *filters.Name != "" {
			if !strings.Contains(strings.ToLower(role.Name), strings.ToLower(*filters.Name)) {
				continue
			}
		}

		// Filtrar por Level
		if filters.Level != nil && role.Level != *filters.Level {
			continue
		}

		filteredRoles = append(filteredRoles, role)
	}

	// Convertir entidades a DTOs
	roleDTOs := make([]domain.RoleDTO, len(filteredRoles))
	for i, role := range filteredRoles {
		roleDTOs[i] = entityToRoleDTO(role)
	}

	uc.log.Info().Int("count", len(roleDTOs)).Msg("Roles obtenidos exitosamente")
	return roleDTOs, nil
}

// entityToRoleDTO convierte una entidad Role a RoleDTO
func entityToRoleDTO(role domain.Role) domain.RoleDTO {
	return domain.RoleDTO{
		ID:               role.ID,
		Name:             role.Name,
		Description:      role.Description,
		Level:            role.Level,
		IsSystem:         role.IsSystem,
		ScopeID:          role.ScopeID,
		ScopeName:        role.ScopeName,
		ScopeCode:        role.ScopeCode,
		BusinessTypeID:   role.BusinessTypeID,
		BusinessTypeName: role.BusinessTypeName,
	}
}
