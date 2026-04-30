package repository

import (
	"context"
	"fmt"
)

func (r *Repository) normalizeCancelledCarrierStatus(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
UPDATE shipments
SET carrier_status        = 'Cancelado',
    carrier_status_detail = 'Cancelado por usuario'
WHERE deleted_at IS NULL
  AND status = 'cancelled'
  AND (carrier_status_detail IS NULL OR carrier_status_detail <> 'Cancelado por usuario')
`).Error; err != nil {
		return fmt.Errorf("normalize cancelled carrier_status: %w", err)
	}

	return nil
}
