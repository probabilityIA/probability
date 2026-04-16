package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

// GetByID obtiene una configuración por su ID
func (uc *useCase) GetByID(ctx context.Context, id uint) (*dtos.NotificationConfigResponseDTO, error) {
	entity, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification config by ID")
		return nil, err
	}

	return mappers.ToResponseDTO(entity), nil
}
