package cache

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// UpdateConfigInCache actualiza una config en cache (puede cambiar IDs)
// NUEVA ESTRUCTURA: Usa IDs en lugar de strings
func (c *cacheManager) UpdateConfigInCache(ctx context.Context, oldConfig, newConfig *entities.IntegrationNotificationConfig) error {
	// 1. Si cambió algún ID de ubicación, eliminar de ubicación vieja
	oldKey := buildCacheKey(oldConfig.IntegrationID, oldConfig.NotificationTypeID, oldConfig.NotificationEventTypeID)
	newKey := buildCacheKey(newConfig.IntegrationID, newConfig.NotificationTypeID, newConfig.NotificationEventTypeID)

	if oldKey != newKey {
		// Eliminar de ubicación vieja
		configIDStr := fmt.Sprintf("%d", oldConfig.ID)
		if err := c.redis.HDel(ctx, oldKey, configIDStr); err != nil {
			c.logger.Warn(ctx).
				Err(err).
				Str("old_key", oldKey).
				Uint("config_id", oldConfig.ID).
				Msg("⚠️  Error eliminando config de ubicación vieja")
		}

		c.logger.Info(ctx).
			Str("old_key", oldKey).
			Str("new_key", newKey).
			Uint("config_id", newConfig.ID).
			Msg("🔄 Config movida a nueva ubicación en cache")
	}

	// 2. Cachear en nueva ubicación (o actualizar en misma ubicación)
	return c.CacheConfig(ctx, newConfig)
}
