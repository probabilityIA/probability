package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateOrderCodIncludesShipping(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Order{}); err != nil {
		return fmt.Errorf("add orders.cod_includes_shipping: %w", err)
	}

	setDefault := `ALTER TABLE orders ALTER COLUMN cod_includes_shipping SET DEFAULT true`
	if err := r.db.Conn(ctx).Exec(setDefault).Error; err != nil {
		return fmt.Errorf("default orders.cod_includes_shipping: %w", err)
	}

	backfill := `UPDATE orders SET cod_includes_shipping = false WHERE platform = 'manual' AND cod_includes_shipping = true`
	if err := r.db.Conn(ctx).Exec(backfill).Error; err != nil {
		return fmt.Errorf("backfill orders.cod_includes_shipping: %w", err)
	}

	seedPlatform := `UPDATE integrations
		SET config = jsonb_set(COALESCE(config, '{}'::jsonb), '{cod_includes_shipping}', 'false'::jsonb)
		WHERE integration_type_id = 6
		  AND deleted_at IS NULL
		  AND (config -> 'cod_includes_shipping') IS NULL`
	if err := r.db.Conn(ctx).Exec(seedPlatform).Error; err != nil {
		return fmt.Errorf("seed cod_includes_shipping en integraciones Plataforma: %w", err)
	}

	return nil
}
