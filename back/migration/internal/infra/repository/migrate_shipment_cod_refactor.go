package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateShipmentCodRefactor(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE shipments
    DROP COLUMN IF EXISTS cod_customer_charge,
    DROP COLUMN IF EXISTS cod_applied_margin
`).Error; err != nil {
		return fmt.Errorf("drop old cod columns: %w", err)
	}

	if err := db.Exec(`
ALTER TABLE shipments
    ADD COLUMN IF NOT EXISTS cod_carrier_fee numeric(12,2),
    ADD COLUMN IF NOT EXISTS cod_probability_margin numeric(12,2)
`).Error; err != nil {
		return fmt.Errorf("add new cod columns: %w", err)
	}

	if err := db.Exec(`DROP TABLE IF EXISTS carrier_cod_config`).Error; err != nil {
		return fmt.Errorf("drop carrier_cod_config: %w", err)
	}

	return nil
}
