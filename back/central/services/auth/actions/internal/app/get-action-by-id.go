package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/auth/actions/internal/domain"
)

// GetActionByID obtiene un action por su ID
func (uc *ActionUseCase) GetActionByID(ctx context.Context, id uint) (*domain.ActionDTO, error) {
	uc.logger.Info().Uint("id", id).Msg("Iniciando obtención de action por ID")

	action, err := uc.repository.GetActionByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error al obtener action")
		return nil, fmt.Errorf("error al obtener action: %w", err)
	}

	if action == nil {
		uc.logger.Warn().Uint("id", id).Msg("Action no encontrado")
		return nil, fmt.Errorf("action con ID %d no encontrado", id)
	}

	actionDTO := &domain.ActionDTO{
		ID:          action.ID,
		Name:        action.Name,
		Description: action.Description,
		CreatedAt:   action.CreatedAt,
		UpdatedAt:   action.UpdatedAt,
	}

	uc.logger.Info().Uint("id", id).Str("name", action.Name).Msg("Action obtenido exitosamente")
	return actionDTO, nil
}
