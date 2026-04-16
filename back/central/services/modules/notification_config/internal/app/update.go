package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

// Update actualiza una configuración de notificación existente
// TODO: Migrar a nueva estructura con IDs (actualmente usando estructura legacy)
func (uc *useCase) Update(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
	// TEMPORALMENTE DESHABILITADO - Necesita migración a nueva estructura
	uc.logger.Warn().
		Uint("id", id).
		Msg("Update endpoint temporarily disabled during migration to new structure")

	// Por ahora, solo permitir actualizar descripción y enabled
	oldConfig, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification config for update")
		return nil, err
	}

	oldConfigCopy := *oldConfig

	// Solo actualizar campos simples que existen en ambas estructuras
	if dto.Description != nil {
		oldConfig.Description = *dto.Description
	}
	if dto.IsActive != nil {
		oldConfig.Enabled = *dto.IsActive // Mapear IsActive a Enabled
	}

	// Persistir cambios
	if err := uc.repository.Update(ctx, oldConfig); err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error updating notification config")
		return nil, err
	}

	// Actualizar en cache
	if err := uc.cacheManager.UpdateConfigInCache(ctx, &oldConfigCopy, oldConfig); err != nil {
		uc.logger.Error().
			Err(err).
			Uint("config_id", id).
			Msg("Error actualizando config en cache")
	}

	updated, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mappers.ToResponseDTO(updated), nil
}
