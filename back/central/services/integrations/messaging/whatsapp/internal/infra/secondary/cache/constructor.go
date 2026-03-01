package cache

import (
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New es un alias para NewNotificationConfigCache para mantener consistencia
func New(redis redisclient.IRedis, logger log.ILogger) ports.INotificationConfigCache {
	return NewNotificationConfigCache(redis, logger)
}
