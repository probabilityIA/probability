package repository

import (
	"context"
	"fmt"
)

func (r *Repository) recomputeShipmentStatusFromWebhook(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
WITH latest AS (
    SELECT DISTINCT ON (tracking_number)
        tracking_number,
        lower(request_body->>'carrier') AS carrier,
        request_body->'events'->-1->>'statusDetail' AS detail,
        (request_body->'events'->-1->>'incidence')::bool AS incidence
    FROM webhook_logs
    WHERE source='envioclick'
      AND tracking_number IS NOT NULL
      AND request_body ? 'events'
      AND jsonb_array_length(request_body->'events') > 0
    ORDER BY tracking_number, created_at DESC
),
mapped AS (
    SELECT tracking_number,
        CASE
            WHEN incidence THEN 'on_hold'
            WHEN carrier='interrapidisimo' AND detail ILIKE 'Envío Admitido' THEN 'picked_up'
            WHEN carrier='servientrega'    AND detail ILIKE 'EN ALISTAMIENTO DEL CLIENTE' THEN 'pending'
            ELSE NULL
        END AS new_status
    FROM latest
)
UPDATE shipments s
SET status = m.new_status
FROM mapped m
WHERE s.tracking_number = m.tracking_number
  AND m.new_status IS NOT NULL
  AND s.deleted_at IS NULL
  AND s.status NOT IN ('cancelled','delivered','returned','failed')
  AND s.status <> m.new_status
`).Error; err != nil {
		return fmt.Errorf("recompute shipment status: %w", err)
	}

	if err := db.Exec(`
UPDATE orders o
SET status = s.status
FROM shipments s
WHERE s.order_id::text = o.id::text
  AND s.deleted_at IS NULL
  AND o.deleted_at IS NULL
  AND s.status IN ('picked_up','pending','on_hold','out_for_delivery')
  AND o.status <> s.status
  AND o.status NOT IN ('cancelled','delivered','returned','failed','refunded')
`).Error; err != nil {
		return fmt.Errorf("propagate to orders: %w", err)
	}

	return nil
}
