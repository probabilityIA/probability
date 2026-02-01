package cache

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// RemoveConfigFromCache elimina una config del cache después de borrarla de BD
func (c *cacheManager) RemoveConfigFromCache(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
	// 1. Obtener todas las keys donde aparece esta config (índice inverso)
	indexKey := buildIndexKey(config.ID)
	keys, err := c.redis.HGetAll(ctx, indexKey)
	if err != nil {
		c.logger.Warn(ctx).
			Err(err).
			Uint("config_id", config.ID).
			Msg("⚠️  Error obteniendo índice inverso - usando key conocida")

		// Continuar con key conocida
		keys = map[string]string{
			buildCacheKey(config.IntegrationID, config.Conditions.Trigger): "1",
		}
	}

	// 2. Eliminar de todas las keys
	configIDStr := fmt.Sprintf("%d", config.ID)
	deletedCount := 0
	for key := range keys {
		if err := c.redis.HDel(ctx, key, configIDStr); err != nil {
			c.logger.Error(ctx).
				Err(err).
				Str("key", key).
				Uint("config_id", config.ID).
				Msg("❌ Error eliminando config de cache")
		} else {
			deletedCount++
		}
	}

	// 3. Limpiar índice inverso
	if err := c.redis.Delete(ctx, indexKey); err != nil {
		c.logger.Warn(ctx).
			Err(err).
			Uint("config_id", config.ID).
			Msg("⚠️  Error limpiando índice inverso")
	}

	c.logger.Info(ctx).
		Uint("config_id", config.ID).
		Int("keys_deleted", deletedCount).
		Msg("✅ Config eliminada del cache exitosamente")

	return nil
}
