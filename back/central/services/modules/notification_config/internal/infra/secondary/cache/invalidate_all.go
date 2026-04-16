package cache

import (
	"context"
	"fmt"
)

// InvalidateAll invalida todo el cache de notification configs
// Limpia primary keys, secondary keys (evt), e índices inversos
func (c *cacheManager) InvalidateAll(ctx context.Context) error {
	// Primary keys: notification:configs:*
	pattern := "notification:configs:*"
	keys, err := c.redis.Keys(ctx, pattern)
	if err != nil {
		c.logger.Error(ctx).Err(err).Msg("❌ Error obteniendo keys para invalidar")
		return fmt.Errorf("error obteniendo keys: %w", err)
	}

	// Eliminar primary + secondary keys (pattern notification:configs:* catches both
	// notification:configs:{id}:{type}:{evt} AND notification:configs:evt:{id}:{code})
	if len(keys) > 0 {
		if err := c.redis.Delete(ctx, keys...); err != nil {
			c.logger.Error(ctx).Err(err).Msg("❌ Error eliminando keys")
			return fmt.Errorf("error eliminando keys: %w", err)
		}
	}

	// Eliminar índices inversos
	indexPattern := "notification:config:*:keys"
	indexKeys, err := c.redis.Keys(ctx, indexPattern)
	if err == nil && len(indexKeys) > 0 {
		if err := c.redis.Delete(ctx, indexKeys...); err != nil {
			c.logger.Warn(ctx).Err(err).Msg("⚠️  Error eliminando índices inversos")
		}
	}

	c.logger.Info(ctx).
		Int("config_keys_deleted", len(keys)).
		Int("index_keys_deleted", len(indexKeys)).
		Msg("✅ Cache completamente invalidado")

	return nil
}
