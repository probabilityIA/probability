package repository

import (
	"context"
	"fmt"
)

func (r *Repository) FixShipmentCodCollectAmount(ctx context.Context) error {
	res := r.db.Conn(ctx).Exec(`
UPDATE shipments s
SET cod_collect_amount = g.cod_value
FROM (
	SELECT DISTINCT ON (l.shipment_id) l.shipment_id, (l.request_payload->>'codValue')::numeric AS cod_value
	FROM shipment_sync_logs l
	WHERE l.request_url LIKE '%/api/v2/shipment'
		AND l.status = 'success'
		AND l.deleted_at IS NULL
		AND (l.request_payload->>'codValue') ~ '^[0-9]+(\.[0-9]+)$|^[0-9]+$'
	ORDER BY l.shipment_id, l.created_at DESC
) g
WHERE g.shipment_id = s.id
	AND g.cod_value > 0
	AND s.cod_collect_amount IS NULL
	AND COALESCE(s.tracking_number, '') <> ''
`)
	if res.Error != nil {
		return fmt.Errorf("rellenar shipments.cod_collect_amount: %w", res.Error)
	}
	fmt.Printf("shipments.cod_collect_amount rellenado en %d envios\n", res.RowsAffected)
	return nil
}
