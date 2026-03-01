package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/events/internal/infra/secondary/cache/mappers"
)

// GetActiveConfigsByIntegrationAndTrigger obtiene configuraciones activas desde Redis cache.
// Lee de la secondary key: notification:configs:evt:{integrationID}:{trigger}
func (c *notificationConfigCache) GetActiveConfigsByIntegrationAndTrigger(
	ctx context.Context,
	integrationID uint,
	trigger string,
) ([]entities.CachedNotificationConfig, error) {
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
		return []entities.CachedNotificationConfig{}, nil
	}

	configs := make([]entities.CachedNotificationConfig, 0, len(entries))
	for configIDStr, jsonData := range entries {
		var cached mappers.CachedNotificationConfigJSON
		if err := json.Unmarshal([]byte(jsonData), &cached); err != nil {
			c.logger.Warn().
				Err(err).
				Str("config_id", configIDStr).
				Msg("Error parseando config desde cache")
			continue
		}

		if !cached.Enabled {
			continue
		}

		configs = append(configs, mappers.FromCachedConfig(&cached))
	}

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
