package cache

import (
	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type sessionCache struct {
	redis redis.IRedis
	log   log.ILogger
}

// New crea un nuevo cache de sesiones AI
func New(redisClient redis.IRedis, logger log.ILogger) domain.ISessionCache {
	return &sessionCache{
		redis: redisClient,
		log:   logger,
	}
}
