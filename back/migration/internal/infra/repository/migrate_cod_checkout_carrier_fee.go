package repository

import (
	"context"
	"fmt"
)

func (r *Repository) MigrateCodCheckoutCarrierFee(ctx context.Context) error {
	err := r.db.Conn(ctx).Exec(`
ALTER TABLE orders
ADD COLUMN IF NOT EXISTS cod_checkout_carrier_fee numeric(12,2) NOT NULL DEFAULT 0
`).Error
	if err != nil {
		return fmt.Errorf("agregar orders.cod_checkout_carrier_fee: %w", err)
	}
	return nil
}
