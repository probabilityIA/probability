package cache

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const statsTTL = time.Hour

type StatsCache struct {
	redis redis.IRedis
	log   log.ILogger
}

func New(redisClient redis.IRedis, logger log.ILogger) domain.IStatsCache {
	if redisClient == nil {
		return nil
	}
	return &StatsCache{
		redis: redisClient,
		log:   logger.WithModule("dashboard.cache"),
	}
}
