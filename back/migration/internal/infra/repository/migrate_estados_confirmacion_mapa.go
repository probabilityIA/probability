package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateEstadosConfirmacionMapa(ctx context.Context) error {
	db := r.db.Conn(ctx)

	var eventoBase, eventoMapa uint
	if err := db.Raw(`
SELECT net.id FROM notification_event_types net
JOIN notification_types nt ON nt.id = net.notification_type_id
WHERE nt.code = 'whatsapp' AND net.event_code = 'order.created' LIMIT 1
`).Scan(&eventoBase).Error; err != nil {
		return fmt.Errorf("lookup evento base: %w", err)
	}
	if err := db.Raw(`
SELECT net.id FROM notification_event_types net
JOIN notification_types nt ON nt.id = net.notification_type_id
WHERE nt.code = 'whatsapp' AND net.event_code = 'order.created_with_map' LIMIT 1
`).Scan(&eventoMapa).Error; err != nil {
		return fmt.Errorf("lookup evento con mapa: %w", err)
	}
	if eventoBase == 0 || eventoMapa == 0 {
		return nil
	}

	if err := db.Exec(`
INSERT INTO notification_event_type_allowed_statuses (notification_event_type_id, order_status_id)
SELECT ?, a.order_status_id
FROM notification_event_type_allowed_statuses a
WHERE a.notification_event_type_id = ?
ON CONFLICT DO NOTHING
`, eventoMapa, eventoBase).Error; err != nil {
		return fmt.Errorf("copiar estados permitidos al evento con mapa: %w", err)
	}

	return nil
}
