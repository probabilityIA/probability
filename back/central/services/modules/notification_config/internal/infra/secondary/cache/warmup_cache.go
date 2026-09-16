package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/cache/mappers"
)

func (c *cacheManager) WarmupCache(ctx context.Context) error {
	filters := dtos.FilterNotificationConfigDTO{
		Enabled: boolPtr(true),
	}

	configs, err := c.repo.List(ctx, filters)
	if err != nil {
		c.logger.Error(ctx).Err(err).Msg("❌ Error obteniendo configs desde BD")
		return fmt.Errorf("error obteniendo configs desde BD: %w", err)
	}

	allStatusIDs := make(map[uint]bool)
	for _, config := range configs {
		for _, id := range config.OrderStatusIDs {
			allStatusIDs[id] = true
		}
	}

	statusCodeMap := make(map[uint]string)
	if len(allStatusIDs) > 0 && c.orderStatusQuerier != nil {
		ids := make([]uint, 0, len(allStatusIDs))
		for id := range allStatusIDs {
			ids = append(ids, id)
		}
		var resolveErr error
		statusCodeMap, resolveErr = c.orderStatusQuerier.GetOrderStatusCodesByIDs(ctx, ids)
		if resolveErr != nil {
			c.logger.Warn(ctx).
				Err(resolveErr).
				Msg("⚠️  Error resolviendo order status codes en warmup - continuando sin ellos")
		}
	}

	cachedCount := 0
	for i := range configs {
		config := &configs[i]
		cachedConfig := mappers.ToCachedConfig(config)

		if len(config.OrderStatusIDs) > 0 {
			codes := make([]string, 0, len(config.OrderStatusIDs))
			for _, id := range config.OrderStatusIDs {
				if code, ok := statusCodeMap[id]; ok {
					codes = append(codes, code)
				}
			}
			cachedConfig.OrderStatusCodes = codes
		}

		configJSON, err := json.Marshal(cachedConfig)
		if err != nil {
			c.logger.Error(ctx).
				Err(err).
				Uint("config_id", config.ID).
				Msg("❌ Error serializando config")
			continue
		}

		configIDStr := fmt.Sprintf("%d", config.ID)
		primaryKey := buildCacheKey(config.IntegrationID, config.NotificationTypeID, config.NotificationEventTypeID)

		if err := c.redis.HSet(ctx, primaryKey, configIDStr, string(configJSON)); err != nil {
			c.logger.Error(ctx).
				Err(err).
				Str("key", primaryKey).
				Uint("config_id", config.ID).
				Msg("❌ Error cacheando config en Redis (primary)")
			continue
		}

		indexKey := buildIndexKey(config.ID)
		if err := c.redis.HSet(ctx, indexKey, primaryKey, "1"); err != nil {
			c.logger.Warn(ctx).
				Err(err).
				Uint("config_id", config.ID).
				Msg("⚠️  Error actualizando índice inverso (primary)")
		}

		c.cacheBusinessWide(ctx, config.NotificationTypeID, config.BusinessID, cachedConfig.EventCode, config.ID, string(configJSON))

		if cachedConfig.EventCode != "" {
			evtKey := buildEventCodeCacheKey(config.IntegrationID, cachedConfig.EventCode)
			if err := c.redis.HSet(ctx, evtKey, configIDStr, string(configJSON)); err != nil {
				c.logger.Warn(ctx).
					Err(err).
					Str("evt_key", evtKey).
					Uint("config_id", config.ID).
					Msg("⚠️  Error cacheando en Redis (secondary evt key)")
			} else {
				if err := c.redis.HSet(ctx, indexKey, evtKey, "1"); err != nil {
					c.logger.Warn(ctx).
						Err(err).
						Uint("config_id", config.ID).
						Msg("⚠️  Error actualizando índice inverso (secondary)")
				}
			}
		}

		cachedCount++
	}

	c.logger.Info(ctx).
		Int("cached_count", cachedCount).
		Int("total_configs", len(configs)).
		Msg("✅ Warmup de cache completado")

	return nil
}
