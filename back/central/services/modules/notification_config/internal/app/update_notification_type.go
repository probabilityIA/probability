package app

import (
	"context"
	"errors"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

// UpdateNotificationType actualiza un tipo de notificación existente
func (uc *useCase) UpdateNotificationType(ctx context.Context, notificationType *entities.NotificationType) error {
	// Validar que el ID existe
	existing, err := uc.notificationTypeRepo.GetByID(ctx, notificationType.ID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", notificationType.ID).Msg("Notification type not found")
		return domainErrors.ErrNotFound
	}

	// Validar campos requeridos
	if strings.TrimSpace(notificationType.Code) == "" {
		uc.logger.Warn().Msg("Code is required for notification type")
		return domainErrors.ErrInvalidInput
	}

	if strings.TrimSpace(notificationType.Name) == "" {
		uc.logger.Warn().Msg("Name is required for notification type")
		return domainErrors.ErrInvalidInput
	}

	// Si el código cambió, verificar que no esté duplicado
	if notificationType.Code != existing.Code {
		existingByCode, err := uc.notificationTypeRepo.GetByCode(ctx, notificationType.Code)
		if err == nil && existingByCode != nil && existingByCode.ID != notificationType.ID {
			uc.logger.Warn().
				Str("code", notificationType.Code).
				Msg("Another notification type with this code already exists")
			return errors.New("notification type with this code already exists")
		}
	}

	// Actualizar
	if err := uc.notificationTypeRepo.Update(ctx, notificationType); err != nil {
		uc.logger.Error().Err(err).Uint("id", notificationType.ID).Msg("Error updating notification type")
		return err
	}

	uc.logger.Info().
		Uint("id", notificationType.ID).
		Str("code", notificationType.Code).
		Msg("Notification type updated successfully")

	return nil
}
