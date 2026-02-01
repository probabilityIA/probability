package cache

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// cacheManager implementa la interfaz ICacheManager
type cacheManager struct {
	redis  redis.IRedis
	repo   ports.IRepository
	logger log.ILogger
}
