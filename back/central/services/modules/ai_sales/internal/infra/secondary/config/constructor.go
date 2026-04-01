package config

import (
	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

type configProvider struct {
	redis redis.IRedis
	log   log.ILogger
}

// New crea un nuevo config provider que lee configuracion AI de platform_creds en Redis
func New(redisClient redis.IRedis, logger log.ILogger) domain.IConfigProvider {
	return &configProvider{
		redis: redisClient,
		log:   logger,
	}
}
