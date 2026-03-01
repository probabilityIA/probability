package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/cache/mappers"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

type notificationConfigCache struct {
	redis  redisclient.IRedis
	logger log.ILogger
}

// NewNotificationConfigCache crea una nueva instancia del cache adapter
func NewNotificationConfigCache(redis redisclient.IRedis, logger log.ILogger) ports.INotificationConfigCache {
	return &notificationConfigCache{
		redis:  redis,
		logger: logger.WithModule("whatsapp-notification-config-cache"),
	}
}

// GetActiveConfigsByIntegrationAndTrigger obtiene configuraciones activas desde el cache de notification_config
// Lee de la secondary key: notification:configs:evt:{integrationID}:{trigger}
func (c *notificationConfigCache) GetActiveConfigsByIntegrationAndTrigger(
	ctx context.Context,
	integrationID uint,
	trigger string,
) ([]dtos.NotificationConfigData, error) {
	// Leer del secondary cache key del módulo notification_config
	key := fmt.Sprintf("notification:configs:evt:%d:%s", integrationID, trigger)

	entries, err := c.redis.HGetAll(ctx, key)
	if err != nil {
		c.logger.Error().
			Err(err).
			Str("key", key).
			Msg("Error obteniendo configs desde Redis cache")
		return nil, fmt.Errorf("error getting configs from Redis: %w", err)
	}

	if len(entries) == 0 {
		c.logger.Debug().
			Str("key", key).
			Uint("integration_id", integrationID).
			Str("trigger", trigger).
			Msg("No hay configs cacheadas para este trigger")
		return []dtos.NotificationConfigData{}, nil
	}

	// Parsear cada entrada JSON → CachedNotificationConfig → NotificationConfigData
	configs := make([]dtos.NotificationConfigData, 0, len(entries))
	for configIDStr, jsonData := range entries {
		var cached mappers.CachedNotificationConfig
		if err := json.Unmarshal([]byte(jsonData), &cached); err != nil {
			c.logger.Warn().
				Err(err).
				Str("config_id", configIDStr).
				Msg("Error parseando config desde cache")
			continue
		}

		// Filtrar por Enabled
		if !cached.Enabled {
			continue
		}

		configs = append(configs, mappers.FromCachedConfig(&cached))
	}

	// Ordenar por ID ascendente (más antiguo primero como prioridad implícita)
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].ID < configs[j].ID
	})

	c.logger.Info().
		Uint("integration_id", integrationID).
		Str("trigger", trigger).
		Int("count", len(configs)).
		Msg("Configs obtenidas desde notification_config secondary cache")

	return configs, nil
}

// ValidateConditions valida si una orden cumple las condiciones de una configuración
func (c *notificationConfigCache) ValidateConditions(
	config *dtos.NotificationConfigData,
	orderStatus string,
	paymentMethodID uint,
	sourceIntegrationID uint,
) bool {
	// 1. Validar statuses (OrderStatusCodes resueltos desde cache)
	if len(config.Statuses) > 0 {
		statusMatch := false
		for _, status := range config.Statuses {
			if status == orderStatus {
				statusMatch = true
				break
			}
		}
		if !statusMatch {
			c.logger.Debug().
				Str("order_status", orderStatus).
				Interface("expected_statuses", config.Statuses).
				Msg("Config no aplica - order_status no está en lista")
			return false
		}
	}

	// 2. Validar payment methods
	if len(config.PaymentMethods) > 0 {
		pmMatch := false
		for _, methodID := range config.PaymentMethods {
			if methodID == paymentMethodID {
				pmMatch = true
				break
			}
		}
		if !pmMatch {
			c.logger.Debug().
				Uint("payment_method_id", paymentMethodID).
				Interface("expected_methods", config.PaymentMethods).
				Msg("Config no aplica - payment_method_id no está en lista")
			return false
		}
	}

	return true
}
