package app

import (
	"context"
	"errors"
	"fmt"

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
	// NUEVA ESTRUCTURA: Buscar por NotificationEventTypeID
	notificationEventTypeID := id
	filters := dtos.FilterNotificationConfigDTO{
		NotificationEventTypeID: &notificationEventTypeID,
	}
	configs, err := uc.repository.List(ctx, filters)
	if err == nil && len(configs) > 0 {
		// Verificar si alguna configuración está activa
		activeConfigs := []string{}
		for _, config := range configs {
			if config.Enabled {
				// Construir descripción de la configuración
				configDesc := ""

				// Usar descripción si existe
				if config.Description != "" {
					configDesc = config.Description
				} else {
					// Si no hay descripción, usar el tipo de notificación si está disponible
					if config.NotificationType != nil {
						configDesc = fmt.Sprintf("Configuración #%d (%s)", config.ID, config.NotificationType.Name)
					} else {
						configDesc = fmt.Sprintf("Configuración #%d", config.ID)
					}
				}

				activeConfigs = append(activeConfigs, configDesc)
			}
		}

		if len(activeConfigs) > 0 {
			uc.logger.Warn().
				Uint("id", id).
				Str("event_code", eventType.EventCode).
				Int("active_configs_count", len(activeConfigs)).
				Strs("active_configs", activeConfigs).
				Msg("Cannot delete notification event type: has active configurations using it")

			// Construir mensaje de error descriptivo
			configList := ""
			for i, cfg := range activeConfigs {
				if i > 0 {
					configList += ", "
				}
				configList += cfg
				// Limitar a 3 ejemplos para no hacer el mensaje muy largo
				if i >= 2 && len(activeConfigs) > 3 {
					configList += "..."
					break
				}
			}
			return errors.New("no se puede eliminar el evento porque está siendo usado por " +
				fmt.Sprintf("%d configuración(es) activa(s): %s. Desactiva o elimina estas configuraciones primero",
				len(activeConfigs), configList))
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
