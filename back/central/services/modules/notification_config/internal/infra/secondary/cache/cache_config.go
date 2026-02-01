package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/cache/mappers"
)

// CacheConfig cachea UNA configuración después de crearla en BD
func (c *cacheManager) CacheConfig(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
	key := buildCacheKey(config.IntegrationID, config.Conditions.Trigger)

	// Serializar a JSON
	cachedConfig := mappers.ToCachedConfig(config)
	configJSON, err := json.Marshal(cachedConfig)
	if err != nil {
		c.logger.Error(ctx).
			Err(err).
			Uint("config_id", config.ID).
			Msg("❌ Error serializando config")
		return fmt.Errorf("error serializando config: %w", err)
	}

	// HSET en Redis
	configIDStr := fmt.Sprintf("%d", config.ID)
	if err := c.redis.HSet(ctx, key, configIDStr, string(configJSON)); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("key", key).
			Uint("config_id", config.ID).
			Msg("❌ Error cacheando en Redis")
		return fmt.Errorf("error cacheando en Redis: %w", err)
	}

	// Actualizar índice inverso
	indexKey := buildIndexKey(config.ID)
	if err := c.redis.HSet(ctx, indexKey, key, "1"); err != nil {
		c.logger.Warn(ctx).
			Err(err).
			Uint("config_id", config.ID).
			Msg("⚠️  Error actualizando índice inverso")
	}

	c.logger.Info(ctx).
		Uint("config_id", config.ID).
		Str("cache_key", key).
		Msg("✅ Config cacheada exitosamente")

	return nil
}
