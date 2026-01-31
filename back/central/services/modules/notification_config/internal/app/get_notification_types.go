package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// GetNotificationTypes obtiene todos los tipos de notificaciones
func (uc *useCase) GetNotificationTypes(ctx context.Context) ([]entities.NotificationType, error) {
	types, err := uc.notificationTypeRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error getting notification types")
		return nil, err
	}

	uc.logger.Info().Int("count", len(types)).Msg("Notification types retrieved successfully")
	return types, nil
}

// GetNotificationTypeByID obtiene un tipo de notificación por su ID
func (uc *useCase) GetNotificationTypeByID(ctx context.Context, id uint) (*entities.NotificationType, error) {
	notifType, err := uc.notificationTypeRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification type by ID")
		return nil, err
	}

	return notifType, nil
}

// GetNotificationTypeByCode obtiene un tipo de notificación por su código
func (uc *useCase) GetNotificationTypeByCode(ctx context.Context, code string) (*entities.NotificationType, error) {
	notifType, err := uc.notificationTypeRepo.GetByCode(ctx, code)
	if err != nil {
		uc.logger.Error().Err(err).Str("code", code).Msg("Error getting notification type by code")
		return nil, err
	}

	return notifType, nil
}
