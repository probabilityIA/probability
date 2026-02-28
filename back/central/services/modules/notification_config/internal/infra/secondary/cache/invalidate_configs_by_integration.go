package cache

import (
	"context"
	"fmt"
)

// InvalidateConfigsByIntegration invalida todas las configs de una integración
// Limpia tanto primary keys como secondary keys (evt)
func (c *cacheManager) InvalidateConfigsByIntegration(ctx context.Context, integrationID uint) error {
	// Primary keys: notification:configs:{integration_id}:*
	primaryPattern := fmt.Sprintf("notification:configs:%d:*", integrationID)
	primaryKeys, err := c.redis.Keys(ctx, primaryPattern)
	if err != nil {
		c.logger.Error(ctx).
			Err(err).
			Uint("integration_id", integrationID).
			Msg("❌ Error obteniendo primary keys para invalidar")
		return fmt.Errorf("error obteniendo keys: %w", err)
	}

	// Secondary keys: notification:configs:evt:{integration_id}:*
	evtPattern := fmt.Sprintf("notification:configs:evt:%d:*", integrationID)
	evtKeys, err := c.redis.Keys(ctx, evtPattern)
	if err != nil {
		c.logger.Warn(ctx).
			Err(err).
			Uint("integration_id", integrationID).
			Msg("⚠️  Error obteniendo evt keys para invalidar")
	}

	// Combinar y eliminar
	allKeys := append(primaryKeys, evtKeys...)
	if len(allKeys) > 0 {
		if err := c.redis.Delete(ctx, allKeys...); err != nil {
			c.logger.Error(ctx).
				Err(err).
				Uint("integration_id", integrationID).
				Msg("❌ Error eliminando keys")
			return fmt.Errorf("error eliminando keys: %w", err)
		}
	}

	c.logger.Info(ctx).
		Uint("integration_id", integrationID).
		Int("primary_keys_deleted", len(primaryKeys)).
		Int("evt_keys_deleted", len(evtKeys)).
		Msg("✅ Cache invalidado para integración")

	return nil
}
