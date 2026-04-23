package cache

import (
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

type ICredentialsCacheMutable interface {
	ports.ICredentialsCache
	SetResolver(resolver ports.IPlatformCredentialsGetter)
}

func New(redis redisclient.IRedis, logger log.ILogger) (ports.IConversationCache, ICredentialsCacheMutable) {
	convCache := newConversationCache(redis, logger)
	credsCache := newCredentialsCache(logger)
	return convCache, credsCache
}
