package cache

import (
	"context"
	"fmt"
)

// InvalidateConfigsByIntegration invalida todas las configs de una integración
func (c *cacheManager) InvalidateConfigsByIntegration(ctx context.Context, integrationID uint) error {
	// Obtener todas las keys que coincidan con el patrón
	pattern := fmt.Sprintf("notification:configs:%d:*", integrationID)
	keys, err := c.redis.Keys(ctx, pattern)
	if err != nil {
		c.logger.Error(ctx).
			Err(err).
			Uint("integration_id", integrationID).
			Msg("❌ Error obteniendo keys para invalidar")
		return fmt.Errorf("error obteniendo keys: %w", err)
	}

	// Eliminar todas las keys
	if len(keys) > 0 {
		if err := c.redis.Delete(ctx, keys...); err != nil {
			c.logger.Error(ctx).
				Err(err).
				Uint("integration_id", integrationID).
				Msg("❌ Error eliminando keys")
			return fmt.Errorf("error eliminando keys: %w", err)
		}
	}

	c.logger.Info(ctx).
		Uint("integration_id", integrationID).
		Int("keys_deleted", len(keys)).
		Msg("✅ Cache invalidado para integración")

	return nil
}
