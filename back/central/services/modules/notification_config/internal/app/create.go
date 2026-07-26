package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

// Create crea una nueva configuración de notificación
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas
func (uc *useCase) Create(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
	// Validar que la configuración no esté duplicada
	// Duplicado = mismo business + misma integración + mismo tipo de notificación + mismo evento
	filters := dtos.FilterNotificationConfigDTO{
		BusinessID:              dto.BusinessID,
		IntegrationID:           &dto.IntegrationID,
		NotificationTypeID:      &dto.NotificationTypeID,
		NotificationEventTypeID: &dto.NotificationEventTypeID,
	}

	existing, err := uc.repository.List(ctx, filters)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error checking for duplicate configs")
		return nil, err
	}

	// Si ya existe una configuración con la misma combinación, es duplicado
	if len(existing) > 0 {
		// Log detallado de TODAS las configuraciones existentes que coinciden
		uc.logger.Warn().
			Int("existing_count", len(existing)).
			Uint("filter_integration_id", dto.IntegrationID).
			Uint("filter_notification_type_id", dto.NotificationTypeID).
			Uint("filter_notification_event_type_id", dto.NotificationEventTypeID).
			Msg("🔍 Checking for duplicates - filters applied")

		for i, cfg := range existing {
			uc.logger.Warn().
				Int("index", i).
				Uint("existing_id", cfg.ID).
				Uint("existing_integration_id", cfg.IntegrationID).
				Uint("existing_notification_type_id", cfg.NotificationTypeID).
				Uint("existing_notification_event_type_id", cfg.NotificationEventTypeID).
				Bool("existing_enabled", cfg.Enabled).
				Msg("❌ Found duplicate config")
		}

		return nil, errors.ErrDuplicateConfig
	}

	// Crear entidad de dominio (nueva estructura)
	entity := &entities.IntegrationNotificationConfig{
		BusinessID:              dto.BusinessID,
		IntegrationID:           dto.IntegrationID,
		NotificationTypeID:      dto.NotificationTypeID,
		NotificationEventTypeID: dto.NotificationEventTypeID,
		Enabled:                 dto.Enabled,
		Description:             dto.Description,
		OrderStatusIDs:          dto.OrderStatusIDs,
		CODOnly:                 dto.CODOnly,
	}

	// Persistir
	if err := uc.repository.Create(ctx, entity); err != nil {
		uc.logger.Error().Err(err).Msg("Error creating notification config")
		return nil, err
	}

	// Cachear en Redis
	if err := uc.cacheManager.CacheConfig(ctx, entity); err != nil {
		uc.logger.Error().
			Err(err).
			Uint("config_id", entity.ID).
			Msg("Error cacheando config después de crear - cache puede estar desincronizado")
		// NO fallar - el cache es secundario
	}

	return mappers.ToResponseDTO(entity), nil
}
