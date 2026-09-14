package repository

import (
	"context"
	"fmt"
)

func (r *Repository) MigrateShipmentCodCollectAmount(ctx context.Context) error {
	err := r.db.Conn(ctx).Exec(`
ALTER TABLE shipments
ADD COLUMN IF NOT EXISTS cod_collect_amount numeric(12,2)
`).Error
	if err != nil {
		return fmt.Errorf("agregar shipments.cod_collect_amount: %w", err)
	}
	return nil
}
