package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

// CreateNotificationEventType crea un nuevo tipo de evento de notificación
func (uc *useCase) CreateNotificationEventType(ctx context.Context, eventType *entities.NotificationEventType) error {
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
	_, err := uc.notificationTypeRepo.GetByID(ctx, eventType.NotificationTypeID)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Uint("notification_type_id", eventType.NotificationTypeID).
			Msg("NotificationType not found")
		return domainErrors.ErrNotFound
	}

	// Crear
	if err := uc.notificationEventRepo.Create(ctx, eventType); err != nil {
		uc.logger.Error().Err(err).Msg("Error creating notification event type")
		return err
	}

	uc.logger.Info().
		Uint("id", eventType.ID).
		Str("event_code", eventType.EventCode).
		Uint("notification_type_id", eventType.NotificationTypeID).
		Msg("Notification event type created successfully")

	return nil
}
