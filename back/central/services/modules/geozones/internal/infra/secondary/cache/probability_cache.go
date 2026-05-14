package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/infra/secondary/metrics"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const (
	probabilityKeyPrefix = "geozones:probability:"
	probabilityTTL       = 10 * time.Minute
)

type ProbabilityCache struct {
	rdb redis.IRedis
	log log.ILogger
}

func NewProbabilityCache(rdb redis.IRedis, logger log.ILogger) ports.IProbabilityCache {
	if rdb != nil {
		rdb.RegisterCachePrefix(probabilityKeyPrefix)
	}
	return &ProbabilityCache{rdb: rdb, log: logger}
}

func probabilityKey(businessID uint, orderID string) string {
	return fmt.Sprintf("%s%d:%s", probabilityKeyPrefix, businessID, orderID)
}

func (c *ProbabilityCache) GetByOrder(ctx context.Context, businessID uint, orderID string) ([]dtos.ProbabilityResult, bool) {
	if c.rdb == nil {
		return nil, false
	}
	val, err := c.rdb.Get(ctx, probabilityKey(businessID, orderID))
	if err != nil || val == "" {
		metrics.ProbabilityCacheOperations.WithLabelValues("miss").Inc()
		return nil, false
	}
	var out []dtos.ProbabilityResult
	if err := json.Unmarshal([]byte(val), &out); err != nil {
		metrics.ProbabilityCacheOperations.WithLabelValues("miss").Inc()
		return nil, false
	}
	metrics.ProbabilityCacheOperations.WithLabelValues("hit").Inc()
	return out, true
}

func (c *ProbabilityCache) SetByOrder(ctx context.Context, businessID uint, orderID string, results []dtos.ProbabilityResult) error {
	if c.rdb == nil {
		return nil
	}
	payload, err := json.Marshal(results)
	if err != nil {
		return err
	}
	if err := c.rdb.Set(ctx, probabilityKey(businessID, orderID), string(payload), probabilityTTL); err != nil {
		c.log.Warn(ctx).Err(err).Str("order_id", orderID).Msg("probability cache set failed")
		return err
	}
	metrics.ProbabilityCacheOperations.WithLabelValues("set").Inc()
	return nil
}

func (c *ProbabilityCache) InvalidateOrder(ctx context.Context, businessID uint, orderID string) error {
	if c.rdb == nil {
		return nil
	}
	if err := c.rdb.Delete(ctx, probabilityKey(businessID, orderID)); err != nil {
		return err
	}
	metrics.ProbabilityCacheOperations.WithLabelValues("invalidate").Inc()
	return nil
}
