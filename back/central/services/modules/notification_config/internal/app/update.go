package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

// Update actualiza una configuración de notificación existente
func (uc *useCase) Update(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
	// Obtener configuración existente
	existing, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification config for update")
		return nil, err
	}

	// Aplicar actualizaciones
	if dto.NotificationType != nil {
		existing.NotificationType = *dto.NotificationType
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	if dto.Conditions != nil {
		existing.Conditions = *dto.Conditions
	}
	if dto.Config != nil {
		existing.Config = *dto.Config
	}
	if dto.Description != nil {
		existing.Description = *dto.Description
	}
	if dto.Priority != nil {
		existing.Priority = *dto.Priority
	}

	// Persistir cambios
	if err := uc.repository.Update(ctx, existing); err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error updating notification config")
		return nil, err
	}

	uc.logger.Info().
		Uint("id", id).
		Msg("Notification config updated successfully")

	// Obtener configuración actualizada
	updated, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mappers.ToResponseDTO(updated), nil
}
