package app

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

// DeleteNotificationEventType elimina un tipo de evento de notificación por su ID
func (uc *useCase) DeleteNotificationEventType(ctx context.Context, id uint) error {
	// Verificar que existe
	eventType, err := uc.notificationEventRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Notification event type not found")
		return domainErrors.ErrNotFound
	}

	// Verificar que no haya configuraciones activas usando este evento
	// Buscamos en la tabla de configuraciones si hay alguna referencia a este trigger
	filters := dtos.FilterNotificationConfigDTO{
		Trigger: &eventType.EventCode,
	}
	configs, err := uc.repository.List(ctx, filters)
	if err == nil && len(configs) > 0 {
		// Verificar si alguna configuración está activa
		hasActive := false
		for _, config := range configs {
			if config.IsActive {
				hasActive = true
				break
			}
		}

		if hasActive {
			uc.logger.Warn().
				Uint("id", id).
				Str("event_code", eventType.EventCode).
				Int("active_configs_count", len(configs)).
				Msg("Cannot delete notification event type: has active configurations using it")
			return errors.New("cannot delete notification event type: has active configurations using this event")
		}
	}

	// Eliminar
	if err := uc.notificationEventRepo.Delete(ctx, id); err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error deleting notification event type")
		return err
	}

	uc.logger.Info().
		Uint("id", id).
		Str("event_code", eventType.EventCode).
		Msg("Notification event type deleted successfully")

	return nil
}
