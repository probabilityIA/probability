package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

// Update actualiza una configuración de notificación existente
func (uc *useCase) Update(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
	// Obtener configuración existente (guardar estado anterior para cache)
	oldConfig, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification config for update")
		return nil, err
	}

	// Crear copia del estado anterior para actualización de cache
	oldConfigCopy := *oldConfig

	// Aplicar actualizaciones
	if dto.NotificationType != nil {
		oldConfig.NotificationType = *dto.NotificationType
	}
	if dto.IsActive != nil {
		oldConfig.IsActive = *dto.IsActive
	}
	if dto.Conditions != nil {
		oldConfig.Conditions = *dto.Conditions
	}
	if dto.Config != nil {
		oldConfig.Config = *dto.Config
	}
	if dto.Description != nil {
		oldConfig.Description = *dto.Description
	}
	if dto.Priority != nil {
		oldConfig.Priority = *dto.Priority
	}

	// Persistir cambios
	if err := uc.repository.Update(ctx, oldConfig); err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error updating notification config")
		return nil, err
	}

	uc.logger.Info().
		Uint("id", id).
		Msg("Notification config updated successfully")

	// Actualizar en cache
	if err := uc.cacheManager.UpdateConfigInCache(ctx, &oldConfigCopy, oldConfig); err != nil {
		uc.logger.Error().
			Err(err).
			Uint("config_id", id).
			Msg("Error actualizando config en cache")
		// NO fallar - el cache es secundario
	}

	// Obtener configuración actualizada
	updated, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mappers.ToResponseDTO(updated), nil
}
