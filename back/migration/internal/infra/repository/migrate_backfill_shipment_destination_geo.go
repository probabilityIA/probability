package repository

import (
	"context"
	"fmt"
)

func (r *Repository) backfillShipmentDestinationGeo(ctx context.Context) error {
	res := r.db.Conn(ctx).Exec(`
		UPDATE shipments s
		SET destination_city  = COALESCE(NULLIF(s.destination_city, ''), o.shipping_city),
		    destination_state = COALESCE(NULLIF(s.destination_state, ''), o.shipping_state)
		FROM orders o
		WHERE o.id = s.order_id
		  AND s.deleted_at IS NULL
		  AND (
		      (s.destination_city IS NULL OR s.destination_city = '')
		      OR (s.destination_state IS NULL OR s.destination_state = '')
		  )
		  AND COALESCE(o.shipping_city, '') <> ''
	`)
	if res.Error != nil {
		return fmt.Errorf("failed to backfill shipment destination geo: %w", res.Error)
	}
	return nil
}
