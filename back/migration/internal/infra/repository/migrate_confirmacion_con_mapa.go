package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateConfirmacionConMapa(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
SELECT setval(
  pg_get_serial_sequence('notification_event_types', 'id'),
  GREATEST(COALESCE((SELECT MAX(id) FROM notification_event_types), 1), 1)
)
`).Error; err != nil {
		return fmt.Errorf("resync secuencia de notification_event_types: %w", err)
	}

	var whatsappTypeID uint
	if err := db.Raw(`SELECT id FROM notification_types WHERE code = 'whatsapp' LIMIT 1`).Scan(&whatsappTypeID).Error; err != nil {
		return fmt.Errorf("lookup tipo whatsapp: %w", err)
	}
	if whatsappTypeID == 0 {
		return nil
	}

	if err := db.Exec(`
INSERT INTO notification_event_types (created_at, updated_at, notification_type_id, event_code, event_name, description, is_active)
VALUES (NOW(), NOW(), ?, 'order.created_with_map', 'Confirmacion de pedido con mapa',
        'Variante de la confirmacion contra entrega que incluye un mapa de la direccion para que el cliente valide su ubicacion antes de generar la guia. Excluyente con la confirmacion normal.', true)
ON CONFLICT (notification_type_id, event_code) DO NOTHING
`, whatsappTypeID).Error; err != nil {
		return fmt.Errorf("insert evento de confirmacion con mapa: %w", err)
	}

	return nil
}
