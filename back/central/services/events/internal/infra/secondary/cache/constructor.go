package cache

import (
	"github.com/secamc93/probability/back/central/services/events/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

type notificationConfigCache struct {
	redis  redisclient.IRedis
	logger log.ILogger
}

// New crea una nueva instancia del cache adapter para configuraciones de notificación
func New(redis redisclient.IRedis, logger log.ILogger) ports.INotificationConfigCache {
	return &notificationConfigCache{
		redis:  redis,
		logger: logger,
	}
}
