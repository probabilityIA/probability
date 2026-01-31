package app

import (
	"context"
	"errors"

	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

// DeleteNotificationType elimina un tipo de notificación por su ID
func (uc *useCase) DeleteNotificationType(ctx context.Context, id uint) error {
	// Verificar que existe
	_, err := uc.notificationTypeRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Notification type not found")
		return domainErrors.ErrNotFound
	}

	// Verificar que no haya configuraciones activas usando este tipo
	// Buscamos en la tabla de configuraciones si hay alguna referencia
	eventTypes, err := uc.notificationEventRepo.GetByNotificationType(ctx, id)
	if err == nil && len(eventTypes) > 0 {
		uc.logger.Warn().
			Uint("id", id).
			Int("event_types_count", len(eventTypes)).
			Msg("Cannot delete notification type: has associated event types")
		return errors.New("cannot delete notification type: has associated event types")
	}

	// Eliminar
	if err := uc.notificationTypeRepo.Delete(ctx, id); err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error deleting notification type")
		return err
	}

	uc.logger.Info().Uint("id", id).Msg("Notification type deleted successfully")
	return nil
}
