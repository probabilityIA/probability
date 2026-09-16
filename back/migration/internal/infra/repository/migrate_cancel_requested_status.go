package repository

import (
	"context"
	"fmt"
)

const solicitudCancelacionBody = "Recibimos tu solicitud para cancelar el pedido {{1}}. Tu pedido ya tiene guía de envío, así que nuestro equipo revisará la cancelación y te confirmará pronto."

func (r *Repository) migrateCancelRequestedStatus(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
SELECT setval(pg_get_serial_sequence('order_statuses', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM order_statuses), 1), 1))
`).Error; err != nil {
		return fmt.Errorf("resync secuencia de order_statuses: %w", err)
	}

	if err := db.Exec(`
INSERT INTO order_statuses (created_at, updated_at, code, name, description, category, is_active, color, priority)
VALUES (NOW(), NOW(), 'cancel_requested', ?, ?, 'issue', true, '#F59E0B', 7)
ON CONFLICT (code) DO NOTHING
`, "Cliente solicita cancelar", "El cliente pidió cancelar por WhatsApp y la orden ya tiene guía. Revisa si la transportadora la recogió y cancélala o retómala.").Error; err != nil {
		return fmt.Errorf("insertar estado cancel_requested: %w", err)
	}

	if err := db.Exec(`
INSERT INTO whatsapp_templates (created_at, updated_at, origin, name, language, category, body_text, meta_template_id, status, scope, header_type)
SELECT NOW(), NOW(), 'system', 'solicitud_cancelacion_recibida', 'es', 'UTILITY', ?, '1782177306268808', 'pending', 'order_event', 'TEXT'
WHERE NOT EXISTS (
	SELECT 1 FROM whatsapp_templates WHERE business_id IS NULL AND name = 'solicitud_cancelacion_recibida' AND deleted_at IS NULL
)
`, solicitudCancelacionBody).Error; err != nil {
		return fmt.Errorf("registrar plantilla solicitud_cancelacion_recibida: %w", err)
	}

	return nil
}
