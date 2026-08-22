package cache

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
)

func (c *StatsCache) Get(ctx context.Context, key string) (*domain.DashboardStats, bool) {
	raw, err := c.redis.Get(ctx, key)
	if err != nil || raw == "" {
		return nil, false
	}

	var stats domain.DashboardStats
	if err := json.Unmarshal([]byte(raw), &stats); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("Cache de dashboard corrupto, se recalcula")
		return nil, false
	}

	return &stats, true
}

func (c *StatsCache) Set(ctx context.Context, key string, stats *domain.DashboardStats) {
	if stats == nil {
		return
	}

	data, err := json.Marshal(stats)
	if err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("No se pudo serializar el dashboard para cache")
		return
	}

	if err := c.redis.Set(ctx, key, string(data), statsTTL); err != nil {
		c.log.Warn(ctx).Err(err).Str("key", key).Msg("No se pudo guardar el dashboard en cache")
	}
}
