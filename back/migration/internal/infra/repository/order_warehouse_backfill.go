package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateOrderWarehouseBackfill(ctx context.Context) error {
	res := r.db.Conn(ctx).Exec(`
		UPDATE orders o
		SET warehouse_id = w.id,
			warehouse_name = CASE WHEN COALESCE(o.warehouse_name, '') = '' THEN w.name ELSE o.warehouse_name END
		FROM (
			SELECT o2.id AS order_id,
				COALESCE(
					(SELECT w1.id FROM warehouses w1
						JOIN integrations i ON i.id = o2.integration_id
						WHERE i.config->>'inventory_warehouse_mode' = 'single'
							AND w1.id::text = i.config->>'inventory_single_warehouse_id'
							AND w1.business_id = o2.business_id
							AND w1.is_active AND w1.deleted_at IS NULL
						LIMIT 1),
					(SELECT w2.id FROM warehouses w2
						WHERE w2.business_id = o2.business_id
							AND w2.is_active AND w2.deleted_at IS NULL
						ORDER BY w2.is_default DESC, w2.id ASC
						LIMIT 1)
				) AS warehouse_id
			FROM orders o2
			WHERE o2.warehouse_id IS NULL AND o2.deleted_at IS NULL
		) t
		JOIN warehouses w ON w.id = t.warehouse_id
		WHERE o.id = t.order_id AND o.warehouse_id IS NULL
	`)
	if res.Error != nil {
		return fmt.Errorf("failed to backfill order warehouses: %w", res.Error)
	}
	fmt.Printf("orders backfilled with warehouse: %d\n", res.RowsAffected)
	return nil
}
