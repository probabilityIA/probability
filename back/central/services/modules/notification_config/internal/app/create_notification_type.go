package app

import (
	"context"
	"errors"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

// CreateNotificationType crea un nuevo tipo de notificación
func (uc *useCase) CreateNotificationType(ctx context.Context, notificationType *entities.NotificationType) error {
	// Validar campos requeridos
	if strings.TrimSpace(notificationType.Code) == "" {
		uc.logger.Warn().Msg("Code is required for notification type")
		return domainErrors.ErrInvalidInput
	}

	if strings.TrimSpace(notificationType.Name) == "" {
		uc.logger.Warn().Msg("Name is required for notification type")
		return domainErrors.ErrInvalidInput
	}

	// Intentar crear
	if err := uc.notificationTypeRepo.Create(ctx, notificationType); err != nil {
		// Detectar error de unicidad de código
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			uc.logger.Warn().
				Str("code", notificationType.Code).
				Msg("Notification type with this code already exists")
			return errors.New("notification type with this code already exists")
		}

		uc.logger.Error().Err(err).Msg("Error creating notification type")
		return err
	}

	uc.logger.Info().
		Uint("id", notificationType.ID).
		Str("code", notificationType.Code).
		Str("name", notificationType.Name).
		Msg("Notification type created successfully")

	return nil
}
