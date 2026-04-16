package cache

import (
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New crea las instancias de cache para WhatsApp (conversation + credentials)
func New(redis redisclient.IRedis, logger log.ILogger) (ports.IConversationCache, ports.ICredentialsCache) {
	convCache := newConversationCache(redis, logger)
	credsCache := newCredentialsCache(redis, logger)
	return convCache, credsCache
}
