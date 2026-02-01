package cache

import (
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New es un alias para NewNotificationConfigCache para mantener consistencia
func New(redis redisclient.IRedis, logger log.ILogger) INotificationConfigCache {
	return NewNotificationConfigCache(redis, logger)
}
