package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateShipmentCarrierStatus(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE shipments
    ADD COLUMN IF NOT EXISTS carrier_status varchar(128),
    ADD COLUMN IF NOT EXISTS carrier_status_detail varchar(255)
`).Error; err != nil {
		return fmt.Errorf("add carrier_status columns: %w", err)
	}

	if err := db.Exec(`
CREATE INDEX IF NOT EXISTS idx_shipments_carrier_status ON shipments (carrier_status)
`).Error; err != nil {
		return fmt.Errorf("create idx_shipments_carrier_status: %w", err)
	}

	return nil
}
