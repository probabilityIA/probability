package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateNotificationEventPickingStatus(ctx context.Context) error {
	db := r.db.Conn(ctx)

	var eventID uint
	if err := db.Raw(`
SELECT net.id
FROM notification_event_types net
JOIN notification_types nt ON nt.id = net.notification_type_id
WHERE nt.code = 'whatsapp' AND net.event_code = 'order.created'
LIMIT 1
`).Scan(&eventID).Error; err != nil {
		return fmt.Errorf("lookup whatsapp order.created event: %w", err)
	}
	if eventID == 0 {
		return nil
	}

	var pickingID uint
	if err := db.Raw(`SELECT id FROM order_statuses WHERE code = 'picking' LIMIT 1`).Scan(&pickingID).Error; err != nil {
		return fmt.Errorf("lookup picking status: %w", err)
	}
	if pickingID == 0 {
		return nil
	}

	if err := db.Exec(`
INSERT INTO notification_event_type_allowed_statuses (notification_event_type_id, order_status_id)
VALUES (?, ?)
ON CONFLICT DO NOTHING
`, eventID, pickingID).Error; err != nil {
		return fmt.Errorf("link picking status to confirmation event: %w", err)
	}

	return nil
}
