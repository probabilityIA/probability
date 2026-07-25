package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateOrderFreeShipping(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Order{}); err != nil {
		return fmt.Errorf("add orders.free_shipping: %w", err)
	}

	setDefault := `ALTER TABLE orders ALTER COLUMN free_shipping SET DEFAULT false`
	if err := r.db.Conn(ctx).Exec(setDefault).Error; err != nil {
		return fmt.Errorf("default orders.free_shipping: %w", err)
	}

	byDiscount := `UPDATE orders
		SET free_shipping = true
		WHERE free_shipping = false
		  AND shipping_cost > 0
		  AND shipping_discount >= shipping_cost`
	if err := r.db.Conn(ctx).Exec(byDiscount).Error; err != nil {
		return fmt.Errorf("backfill free_shipping por descuento: %w", err)
	}

	byShippingLine := `UPDATE orders
		SET free_shipping = true
		WHERE free_shipping = false
		  AND shipping_cost = 0
		  AND jsonb_typeof(shipping_details -> 'shipping_lines') = 'array'
		  AND EXISTS (
			SELECT 1
			FROM jsonb_array_elements(shipping_details -> 'shipping_lines') AS line
			WHERE COALESCE(line ->> 'price', line ->> 'total', '') IN ('0', '0.0', '0.00')
		  )`
	if err := r.db.Conn(ctx).Exec(byShippingLine).Error; err != nil {
		return fmt.Errorf("backfill free_shipping por linea de envio: %w", err)
	}

	return nil
}
