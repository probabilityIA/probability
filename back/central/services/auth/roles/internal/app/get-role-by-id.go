package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
)

// GetRoleByID obtiene un rol por su ID
func (uc *RoleUseCase) GetRoleByID(ctx context.Context, id uint) (*domain.RoleDTO, error) {
	uc.log.Info().Uint("id", id).Msg("Iniciando caso de uso: obtener rol por ID")

	role, err := uc.repository.GetRoleByID(ctx, id)
	if err != nil {
		uc.log.Error().Uint("id", id).Err(err).Msg("Error al obtener rol por ID desde el repositorio")
		return nil, fmt.Errorf("rol no encontrado")
	}

	if role == nil {
		uc.log.Error().Uint("id", id).Msg("Rol no encontrado")
		return nil, fmt.Errorf("rol no encontrado")
	}

	roleDTO := entityToRoleDTO(*role)

	uc.log.Info().Uint("id", id).Msg("Rol obtenido exitosamente")
	return &roleDTO, nil
}
