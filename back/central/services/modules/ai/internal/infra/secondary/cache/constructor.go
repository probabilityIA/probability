package cache

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const (
	usagePrefix = "ai:assistant:usage:v1"
	introPrefix = "ai:assistant:intro:v1"
	introTTL    = 365 * 24 * time.Hour
)

type Store struct {
	redis redis.IRedis
}

func New(client redis.IRedis) ports.IAssistantStore {
	if client == nil {
		return nil
	}
	client.RegisterCachePrefix(usagePrefix)
	client.RegisterCachePrefix(introPrefix)
	return &Store{redis: client}
}
