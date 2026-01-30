package app

import (
	"context"
)

// Delete elimina una configuración de notificación
func (uc *useCase) Delete(ctx context.Context, id uint) error {
	if err := uc.repository.Delete(ctx, id); err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error deleting notification config")
		return err
	}

	uc.logger.Info().
		Uint("id", id).
		Msg("Notification config deleted successfully")

	return nil
}
