package repository

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateShippingConfigs(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.ShippingConfig{}); err != nil {
		return err
	}

	if err := r.db.Conn(ctx).Exec(`
		INSERT INTO shipping_configs (
			created_at, updated_at, business_id, warehouse_id,
			package_strategy, boxes, carriers
		)
		SELECT now(), now(), w.business_id, w.id,
			w.metadata->>'shipping_package_strategy',
			coalesce((w.metadata::jsonb)->'standard_boxes', '[]'::jsonb),
			'[]'::jsonb
		FROM warehouses w
		WHERE w.deleted_at IS NULL
		  AND w.metadata->>'shipping_package_strategy' IS NOT NULL
		  AND NOT EXISTS (
			SELECT 1 FROM shipping_configs sc
			WHERE sc.business_id = w.business_id AND sc.warehouse_id = w.id
		  )
	`).Error; err != nil {
		return err
	}

	return nil
}
