package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/cache/mappers"
)

// WarmupCache carga todas las configuraciones activas en Redis al iniciar el sistema
func (c *cacheManager) WarmupCache(ctx context.Context) error {
	c.logger.Info(ctx).Msg("🔄 Iniciando warmup de cache de configuraciones de notificación")

	// 1. Obtener TODAS las configuraciones activas desde BD
	filters := dtos.FilterNotificationConfigDTO{
		IsActive: boolPtr(true),
	}

	configs, err := c.repo.List(ctx, filters)
	if err != nil {
		c.logger.Error(ctx).Err(err).Msg("❌ Error obteniendo configs desde BD")
		return fmt.Errorf("error obteniendo configs desde BD: %w", err)
	}

	c.logger.Info(ctx).Int("count", len(configs)).Msg("📊 Configuraciones obtenidas desde BD")

	// 2. Agrupar por integration_id + trigger
	grouped := make(map[string][]*entities.IntegrationNotificationConfig)
	for i := range configs {
		config := &configs[i]
		key := buildCacheKey(config.IntegrationID, config.Conditions.Trigger)
		grouped[key] = append(grouped[key], config)
	}

	c.logger.Info(ctx).Int("cache_keys", len(grouped)).Msg("🗂️  Configs agrupadas por integration+trigger")

	// 3. Cachear en Redis
	cachedCount := 0
	for key, configList := range grouped {
		for _, config := range configList {
			// Serializar config a JSON
			cachedConfig := mappers.ToCachedConfig(config)
			configJSON, err := json.Marshal(cachedConfig)
			if err != nil {
				c.logger.Error(ctx).
					Err(err).
					Uint("config_id", config.ID).
					Msg("❌ Error serializando config")
				continue
			}

			// HSET: key = notification:configs:{integration}:{trigger}, field = {config_id}, value = JSON
			configIDStr := fmt.Sprintf("%d", config.ID)
			if err := c.redis.HSet(ctx, key, configIDStr, string(configJSON)); err != nil {
				c.logger.Error(ctx).
					Err(err).
					Str("key", key).
					Uint("config_id", config.ID).
					Msg("❌ Error cacheando config en Redis")
				continue
			}

			// Actualizar índice inverso
			indexKey := buildIndexKey(config.ID)
			if err := c.redis.HSet(ctx, indexKey, key, "1"); err != nil {
				c.logger.Warn(ctx).
					Err(err).
					Uint("config_id", config.ID).
					Msg("⚠️  Error actualizando índice inverso")
			}

			cachedCount++
		}
	}

	c.logger.Info(ctx).
		Int("total_configs", len(configs)).
		Int("cache_keys", len(grouped)).
		Int("cached_configs", cachedCount).
		Msg("✅ Cache warmup completado exitosamente")

	return nil
}
