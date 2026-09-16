package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/cache/mappers"
)

func (c *cacheManager) CacheConfig(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
	key := buildCacheKey(config.IntegrationID, config.NotificationTypeID, config.NotificationEventTypeID)

	cachedConfig := mappers.ToCachedConfig(config)

	if len(config.OrderStatusIDs) > 0 && c.orderStatusQuerier != nil {
		codeMap, err := c.orderStatusQuerier.GetOrderStatusCodesByIDs(ctx, config.OrderStatusIDs)
		if err != nil {
			c.logger.Warn(ctx).
				Err(err).
				Uint("config_id", config.ID).
				Msg("⚠️  Error resolviendo order status codes - continuando sin ellos")
		} else {
			codes := make([]string, 0, len(codeMap))
			for _, code := range codeMap {
				codes = append(codes, code)
			}
			cachedConfig.OrderStatusCodes = codes
		}
	}

	if cachedConfig.EventCode == "" && config.NotificationEventTypeID > 0 {
		fullConfig, err := c.repo.GetByID(ctx, config.ID)
		if err == nil && fullConfig.NotificationEventType != nil {
			cachedConfig.EventCode = fullConfig.NotificationEventType.EventCode
		}
	}

	configJSON, err := json.Marshal(cachedConfig)
	if err != nil {
		c.logger.Error(ctx).
			Err(err).
			Uint("config_id", config.ID).
			Msg("❌ Error serializando config")
		return fmt.Errorf("error serializando config: %w", err)
	}

	configIDStr := fmt.Sprintf("%d", config.ID)

	if err := c.redis.HSet(ctx, key, configIDStr, string(configJSON)); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("key", key).
			Uint("config_id", config.ID).
			Msg("❌ Error cacheando en Redis (primary key)")
		return fmt.Errorf("error cacheando en Redis: %w", err)
	}

	indexKey := buildIndexKey(config.ID)
	if err := c.redis.HSet(ctx, indexKey, key, "1"); err != nil {
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

	c.logger.Info(ctx).
		Uint("config_id", config.ID).
		Str("cache_key", key).
		Str("event_code", cachedConfig.EventCode).
		Msg("✅ Config cacheada exitosamente")

	return nil
}
