package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateWalletTxShipmentLink(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE transaction
    ADD COLUMN IF NOT EXISTS shipment_id BIGINT
`).Error; err != nil {
		return fmt.Errorf("add shipment_id column: %w", err)
	}

	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_transaction_shipment_id ON transaction(shipment_id)`).Error; err != nil {
		return fmt.Errorf("create shipment_id index: %w", err)
	}

	if err := db.Exec(`
ALTER TABLE transaction
    DROP CONSTRAINT IF EXISTS fk_transaction_shipment
`).Error; err != nil {
		return fmt.Errorf("drop existing fk: %w", err)
	}

	if err := db.Exec(`
ALTER TABLE transaction
    ADD CONSTRAINT fk_transaction_shipment
    FOREIGN KEY (shipment_id) REFERENCES shipments(id) ON UPDATE CASCADE ON DELETE SET NULL
`).Error; err != nil {
		return fmt.Errorf("add shipment_id fk: %w", err)
	}

	if err := db.Exec(`
UPDATE transaction t
SET shipment_id = s.id
FROM shipments s
WHERE t.shipment_id IS NULL
  AND t.reference LIKE '%Guide generation: %'
  AND s.tracking_number IS NOT NULL
  AND s.deleted_at IS NULL
  AND substring(t.reference from 'Guide generation: (.+)$') = s.tracking_number
`).Error; err != nil {
		return fmt.Errorf("backfill shipment_id from reference: %w", err)
	}

	return nil
}
