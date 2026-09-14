package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const (
	keyPrefix = "authz:access:v1"
	ttl       = 60 * time.Second
)

type Cache struct {
	redis redis.IRedis
}

func New(client redis.IRedis) ports.IAccessCache {
	if client == nil {
		return nil
	}
	client.RegisterCachePrefix(keyPrefix)
	return &Cache{redis: client}
}

func key(userID, businessID uint) string {
	return fmt.Sprintf("%s:%d:%d", keyPrefix, userID, businessID)
}

func (c *Cache) Get(ctx context.Context, userID, businessID uint) (*entities.Access, bool) {
	raw, err := c.redis.Get(ctx, key(userID, businessID))
	if err != nil || raw == "" {
		return nil, false
	}
	var access entities.Access
	if err := json.Unmarshal([]byte(raw), &access); err != nil {
		return nil, false
	}
	return &access, true
}

func (c *Cache) Set(ctx context.Context, access *entities.Access) {
	raw, err := json.Marshal(access)
	if err != nil {
		return
	}
	_ = c.redis.Set(ctx, key(access.UserID, access.BusinessID), string(raw), ttl)
}
