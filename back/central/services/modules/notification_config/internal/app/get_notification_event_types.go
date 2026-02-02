package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// GetEventTypesByNotificationType obtiene todos los tipos de eventos de un tipo de notificación
func (uc *useCase) GetEventTypesByNotificationType(ctx context.Context, notificationTypeID uint) ([]entities.NotificationEventType, error) {
	eventTypes, err := uc.notificationEventRepo.GetByNotificationType(ctx, notificationTypeID)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Uint("notification_type_id", notificationTypeID).
			Msg("Error getting event types by notification type")
		return nil, err
	}

	uc.logger.Info().
		Uint("notification_type_id", notificationTypeID).
		Int("count", len(eventTypes)).
		Msg("Event types retrieved successfully")

	return eventTypes, nil
}

// GetNotificationEventTypeByID obtiene un tipo de evento de notificación por su ID
func (uc *useCase) GetNotificationEventTypeByID(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
	eventType, err := uc.notificationEventRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification event type by ID")
		return nil, err
	}

	return eventType, nil
}

// ListAllEventTypes obtiene todos los tipos de eventos de notificación
func (uc *useCase) ListAllEventTypes(ctx context.Context) ([]entities.NotificationEventType, error) {
	uc.logger.Info().Msg("🔍 Fetching all notification event types from repository")

	eventTypes, err := uc.notificationEventRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error().Err(err).Msg("❌ Error getting all event types")
		return nil, err
	}

	uc.logger.Info().Int("count", len(eventTypes)).Msg("✅ All event types retrieved successfully")

	return eventTypes, nil
}
