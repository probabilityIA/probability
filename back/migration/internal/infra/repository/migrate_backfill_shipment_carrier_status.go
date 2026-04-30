package repository

import (
	"context"
	"fmt"
)

func (r *Repository) backfillShipmentCarrierStatus(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
WITH latest_per_tracking AS (
    SELECT DISTINCT ON (tracking_number)
        tracking_number,
        request_body
    FROM webhook_logs
    WHERE source = 'envioclick'
      AND tracking_number IS NOT NULL
      AND request_body ? 'events'
      AND jsonb_array_length(request_body->'events') > 0
    ORDER BY tracking_number, created_at DESC
)
UPDATE shipments s
SET
    carrier_status        = NULLIF(l.request_body->'events'->-1->>'statusStep', ''),
    carrier_status_detail = NULLIF(l.request_body->'events'->-1->>'statusDetail', '')
FROM latest_per_tracking l
WHERE s.tracking_number = l.tracking_number
  AND s.deleted_at IS NULL
`).Error; err != nil {
		return fmt.Errorf("backfill carrier_status from webhook_logs: %w", err)
	}

	return nil
}
