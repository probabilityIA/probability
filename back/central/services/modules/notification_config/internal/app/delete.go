package app

import (
	"context"
)

// Delete elimina una configuración de notificación
func (uc *useCase) Delete(ctx context.Context, id uint) error {
	// Obtener config antes de eliminar (necesario para remover del cache)
	config, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error getting notification config for delete")
		return err
	}

	// Eliminar de BD
	if err := uc.repository.Delete(ctx, id); err != nil {
		uc.logger.Error().Err(err).Uint("id", id).Msg("Error deleting notification config")
		return err
	}

	uc.logger.Info().
		Uint("id", id).
		Msg("Notification config deleted successfully")

	// Eliminar de cache
	if err := uc.cacheManager.RemoveConfigFromCache(ctx, config); err != nil {
		uc.logger.Error().
			Err(err).
			Uint("config_id", id).
			Msg("Error eliminando config del cache")
		// NO fallar - el cache es secundario
	}

	return nil
}
