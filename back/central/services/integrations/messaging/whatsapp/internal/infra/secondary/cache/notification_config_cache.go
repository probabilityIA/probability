package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/cache/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/repository"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/log"
)

// INotificationConfigCache define la interfaz para consultar configuraciones desde Redis (read-only)
type INotificationConfigCache interface {
	GetActiveConfigsByIntegrationAndTrigger(ctx context.Context, integrationID uint, trigger string) ([]repository.NotificationConfigData, error)
	ValidateConditions(config *repository.NotificationConfigData, orderStatus string, paymentMethodID uint, sourceIntegrationID uint) bool
}

type notificationConfigCache struct {
	redis  redisclient.IRedis
	logger log.ILogger
}

// NewNotificationConfigCache crea una nueva instancia del cache adapter
func NewNotificationConfigCache(redis redisclient.IRedis, logger log.ILogger) INotificationConfigCache {
	return &notificationConfigCache{
		redis:  redis,
		logger: logger.WithModule("whatsapp-notification-config-cache"),
	}
}

// GetActiveConfigsByIntegrationAndTrigger obtiene configuraciones activas desde Redis cache
func (c *notificationConfigCache) GetActiveConfigsByIntegrationAndTrigger(
	ctx context.Context,
	integrationID uint,
	trigger string,
) ([]repository.NotificationConfigData, error) {
	// Construir key de Redis
	key := fmt.Sprintf("notification:configs:%d:%s", integrationID, trigger)

	// HGETALL - obtiene todos los fields del hash
	results, err := c.redis.HGetAll(ctx, key)
	if err != nil {
		c.logger.Error().
			Err(err).
			Str("key", key).
			Msg("❌ Error obteniendo configs desde Redis")
		return nil, fmt.Errorf("error getting configs from Redis: %w", err)
	}

	// Si no hay resultados, retornar array vacío (no es error)
	if len(results) == 0 {
		c.logger.Debug().
			Str("key", key).
			Msg("ℹ️  No hay configuraciones cacheadas para esta integración y trigger")
		return []repository.NotificationConfigData{}, nil
	}

	// Deserializar cada config
	configs := make([]repository.NotificationConfigData, 0, len(results))
	for field, jsonData := range results {
		var cachedConfig mappers.CachedNotificationConfig
		if err := json.Unmarshal([]byte(jsonData), &cachedConfig); err != nil {
			c.logger.Error().
				Err(err).
				Str("field", field).
				Str("key", key).
				Msg("❌ Error deserializando config desde cache - saltando")
			continue
		}

		// Solo agregar configs activas
		if cachedConfig.IsActive {
			configs = append(configs, mappers.FromCachedConfig(&cachedConfig))
		}
	}

	// Ordenar por prioridad descendente
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Priority > configs[j].Priority
	})

	c.logger.Info().
		Str("key", key).
		Int("count", len(configs)).
		Msg("✅ Configuraciones obtenidas desde Redis cache")

	return configs, nil
}

// ValidateConditions valida si una orden cumple las condiciones de una configuración
func (c *notificationConfigCache) ValidateConditions(
	config *repository.NotificationConfigData,
	orderStatus string,
	paymentMethodID uint,
	sourceIntegrationID uint,
) bool {
	// 1. Validar source_integration_id PRIMERO (más específico)
	if config.SourceIntegrationID != nil {
		// Si la config especifica una integración origen, DEBE coincidir
		if *config.SourceIntegrationID != sourceIntegrationID {
			c.logger.Debug().
				Uint("expected", *config.SourceIntegrationID).
				Uint("actual", sourceIntegrationID).
				Msg("Config no aplica - source_integration_id no coincide")
			return false
		}
	}
	// Si config.SourceIntegrationID == nil → aplica a todas las integraciones

	// 2. Validar statuses
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

	// 3. Validar payment methods
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

	// Todas las validaciones pasaron
	return true
}
