package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

const ttlStats = time.Minute

func statsKey(businessID uint) string {
	return fmt.Sprintf("integration:stats:biz:%d", businessID)
}

func (c *IntegrationCache) SetIntegrationStats(ctx context.Context, businessID uint, stats []domain.IntegrationStats) error {
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, statsKey(businessID), string(data), ttlStats)
}

func (c *IntegrationCache) GetIntegrationStats(ctx context.Context, businessID uint) ([]domain.IntegrationStats, error) {
	data, err := c.redis.Get(ctx, statsKey(businessID))
	if err != nil {
		return nil, err
	}
	var stats []domain.IntegrationStats
	if err := json.Unmarshal([]byte(data), &stats); err != nil {
		return nil, err
	}
	return stats, nil
}
