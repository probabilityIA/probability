package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

// UpdateNotificationEventType actualiza un tipo de evento de notificación existente
func (uc *useCase) UpdateNotificationEventType(ctx context.Context, eventType *entities.NotificationEventType) error {
	// Verificar que el ID existe
	_, err := uc.notificationEventRepo.GetByID(ctx, eventType.ID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", eventType.ID).Msg("Notification event type not found")
		return domainErrors.ErrNotFound
	}

	// Validar campos requeridos
	if eventType.NotificationTypeID == 0 {
		uc.logger.Warn().Msg("NotificationTypeID is required for notification event type")
		return domainErrors.ErrInvalidInput
	}

	if strings.TrimSpace(eventType.EventCode) == "" {
		uc.logger.Warn().Msg("EventCode is required for notification event type")
		return domainErrors.ErrInvalidInput
	}

	if strings.TrimSpace(eventType.EventName) == "" {
		uc.logger.Warn().Msg("EventName is required for notification event type")
		return domainErrors.ErrInvalidInput
	}

	// Verificar que el NotificationType existe
	_, err = uc.notificationTypeRepo.GetByID(ctx, eventType.NotificationTypeID)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Uint("notification_type_id", eventType.NotificationTypeID).
			Msg("NotificationType not found")
		return domainErrors.ErrNotFound
	}

	// Actualizar
	if err := uc.notificationEventRepo.Update(ctx, eventType); err != nil {
		uc.logger.Error().Err(err).Uint("id", eventType.ID).Msg("Error updating notification event type")
		return err
	}

	uc.logger.Info().
		Uint("id", eventType.ID).
		Str("event_code", eventType.EventCode).
		Msg("Notification event type updated successfully")

	return nil
}
