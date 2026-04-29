package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateShipmentCarrierCost(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE shipments
    ADD COLUMN IF NOT EXISTS carrier_cost numeric(12,2),
    ADD COLUMN IF NOT EXISTS applied_margin numeric(12,2)
`).Error; err != nil {
		return fmt.Errorf("add carrier_cost columns: %w", err)
	}

	if err := db.Exec(`
UPDATE shipments s
SET carrier_cost = GREATEST(s.total_cost - sm.margin_amount, 0),
    applied_margin = sm.margin_amount
FROM (
    SELECT business_id, lower(carrier_code) AS code, MAX(margin_amount) AS margin_amount
    FROM shipping_margin
    WHERE deleted_at IS NULL
    GROUP BY business_id, lower(carrier_code)
) sm,
    orders o
WHERE s.deleted_at IS NULL
  AND s.carrier_cost IS NULL
  AND s.total_cost IS NOT NULL
  AND s.order_id = o.id
  AND o.business_id = sm.business_id
  AND lower(COALESCE(s.carrier, '')) = sm.code
`).Error; err != nil {
		return fmt.Errorf("backfill shipments carrier_cost: %w", err)
	}

	return nil
}
