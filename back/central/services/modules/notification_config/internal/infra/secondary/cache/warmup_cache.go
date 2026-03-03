package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/cache/mappers"
)

// WarmupCache carga todas las configuraciones activas en Redis al iniciar el sistema
// Escribe en keys primarias Y secondary keys (por event_code)
func (c *cacheManager) WarmupCache(ctx context.Context) error {
	// 1. Obtener TODAS las configuraciones activas desde BD
	// List() ya hace Preload("NotificationType") y Preload("NotificationEventType")
	filters := dtos.FilterNotificationConfigDTO{
		Enabled: boolPtr(true),
	}

	configs, err := c.repo.List(ctx, filters)
	if err != nil {
		c.logger.Error(ctx).Err(err).Msg("❌ Error obteniendo configs desde BD")
		return fmt.Errorf("error obteniendo configs desde BD: %w", err)
	}

	// 2. Recolectar todos los OrderStatusIDs únicos para resolver codes en batch
	allStatusIDs := make(map[uint]bool)
	for _, config := range configs {
		for _, id := range config.OrderStatusIDs {
			allStatusIDs[id] = true
		}
	}

	// Resolver OrderStatusCodes en batch
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

	// 3. Agrupar por keys y cachear
	cachedCount := 0
	for i := range configs {
		config := &configs[i]
		cachedConfig := mappers.ToCachedConfig(config)

		// Resolver OrderStatusCodes desde el map batch
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

		// HSET en key primaria
		if err := c.redis.HSet(ctx, primaryKey, configIDStr, string(configJSON)); err != nil {
			c.logger.Error(ctx).
				Err(err).
				Str("key", primaryKey).
				Uint("config_id", config.ID).
				Msg("❌ Error cacheando config en Redis (primary)")
			continue
		}

		// Actualizar índice inverso con primary key
		indexKey := buildIndexKey(config.ID)
		if err := c.redis.HSet(ctx, indexKey, primaryKey, "1"); err != nil {
			c.logger.Warn(ctx).
				Err(err).
				Uint("config_id", config.ID).
				Msg("⚠️  Error actualizando índice inverso (primary)")
		}

		// HSET en key secundaria (por event_code) si tenemos EventCode
		if cachedConfig.EventCode != "" {
			evtKey := buildEventCodeCacheKey(config.IntegrationID, cachedConfig.EventCode)
			if err := c.redis.HSet(ctx, evtKey, configIDStr, string(configJSON)); err != nil {
				c.logger.Warn(ctx).
					Err(err).
					Str("evt_key", evtKey).
					Uint("config_id", config.ID).
					Msg("⚠️  Error cacheando en Redis (secondary evt key)")
			} else {
				// Registrar secondary key en índice inverso
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
