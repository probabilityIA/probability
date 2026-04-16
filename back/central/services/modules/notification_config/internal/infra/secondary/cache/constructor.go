package cache

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// New crea una nueva instancia del cache manager
func New(redis redis.IRedis, repo ports.IRepository, orderStatusQuerier ports.IOrderStatusQuerier, logger log.ILogger) ports.ICacheManager {
	return &cacheManager{
		redis:              redis,
		repo:               repo,
		orderStatusQuerier: orderStatusQuerier,
		logger:             logger.WithModule("notification-config-cache"),
	}
}

