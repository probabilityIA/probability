package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

// List obtiene una lista de configuraciones con filtros
func (uc *useCase) List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error) {
	entities, err := uc.repository.List(ctx, filters)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error listing notification configs")
		return nil, err
	}

	return mappers.ToResponseDTOList(entities), nil
}
