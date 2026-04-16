package cache

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// UpdateConfigInCache actualiza una config en cache (puede cambiar IDs)
// Maneja tanto primary keys como secondary keys (evt)
func (c *cacheManager) UpdateConfigInCache(ctx context.Context, oldConfig, newConfig *entities.IntegrationNotificationConfig) error {
	// 1. Si cambió algún ID de ubicación, eliminar de ubicación vieja
	oldKey := buildCacheKey(oldConfig.IntegrationID, oldConfig.NotificationTypeID, oldConfig.NotificationEventTypeID)
	newKey := buildCacheKey(newConfig.IntegrationID, newConfig.NotificationTypeID, newConfig.NotificationEventTypeID)
	configIDStr := fmt.Sprintf("%d", oldConfig.ID)

	if oldKey != newKey {
		// Eliminar de primary key vieja
		if err := c.redis.HDel(ctx, oldKey, configIDStr); err != nil {
			c.logger.Warn(ctx).
				Err(err).
				Str("old_key", oldKey).
				Uint("config_id", oldConfig.ID).
				Msg("⚠️  Error eliminando config de ubicación vieja (primary)")
		}

		c.logger.Info(ctx).
			Str("old_key", oldKey).
			Str("new_key", newKey).
			Uint("config_id", newConfig.ID).
			Msg("🔄 Config movida a nueva ubicación en cache")
	}

	// 2. Limpiar secondary key vieja si tenemos EventCode del oldConfig
	if oldConfig.NotificationEventType != nil && oldConfig.NotificationEventType.EventCode != "" {
		oldEvtKey := buildEventCodeCacheKey(oldConfig.IntegrationID, oldConfig.NotificationEventType.EventCode)
		if err := c.redis.HDel(ctx, oldEvtKey, configIDStr); err != nil {
			c.logger.Warn(ctx).
				Err(err).
				Str("old_evt_key", oldEvtKey).
				Uint("config_id", oldConfig.ID).
				Msg("⚠️  Error eliminando config de secondary key vieja")
		}
	}

	// 3. Cachear en nueva ubicación (CacheConfig escribe primary + secondary)
	return c.CacheConfig(ctx, newConfig)
}
