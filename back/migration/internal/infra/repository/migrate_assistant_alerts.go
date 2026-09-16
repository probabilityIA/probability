package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const assistantNotificationTypeID = 6

type assistantEventSeed struct {
	code string
	name string
	desc string
}

func (r *Repository) migrateAssistantAlerts(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.AssistantAlert{}, &models.AssistantAlertCursor{}); err != nil {
		return fmt.Errorf("failed to auto-migrate assistant alerts: %w", err)
	}

	if err := db.Exec(`
INSERT INTO notification_types (id, created_at, updated_at, name, code, description, icon, is_active)
VALUES (?, NOW(), NOW(), ?, 'assistant', 'Alertas dentro de la plataforma, en el chat del asistente Via', 'sparkles', true)
ON CONFLICT (code) DO NOTHING
`, assistantNotificationTypeID, "V\u00eda").Error; err != nil {
		return fmt.Errorf("insert assistant notification type: %w", err)
	}

	var typeID uint
	if err := db.Raw(`SELECT id FROM notification_types WHERE code = 'assistant' LIMIT 1`).Scan(&typeID).Error; err != nil {
		return fmt.Errorf("lookup assistant notification type: %w", err)
	}
	if typeID != assistantNotificationTypeID {
		return fmt.Errorf("el canal assistant quedo con id %d y el dispatcher espera %d", typeID, assistantNotificationTypeID)
	}

	for _, table := range []string{"notification_types", "notification_event_types"} {
		if err := db.Exec(`
SELECT setval(
  pg_get_serial_sequence(?, 'id'),
  GREATEST(COALESCE((SELECT MAX(id) FROM `+table+`), 1), 1)
)
`, table).Error; err != nil {
			return fmt.Errorf("resync secuencia de %s: %w", table, err)
		}
	}

	seeds := []assistantEventSeed{
		{"order.cancelled", "El cliente pidio cancelar", "El cliente cancelo el pedido desde WhatsApp"},
		{"order.status_changed", "Cambio de estado de la orden", "Elige que estados avisan: cancelada, fallida, rechazada, en espera, reembolsada, novedad de entrega, entrega fallida, devuelta, novedad de inventario o entregada"},
		{"shipment.guide_failed", "Guia rechazada", "La transportadora no genero la guia"},
		{"shipment.tracking_updated", "Novedad de transportadora", "Novedad, devolucion o entrega fallida reportada por la transportadora"},
		{"shipment.cancel_failed", "Cancelacion de guia fallida", "La transportadora no acepto cancelar la guia"},
		{"invoice.failed", "Factura rechazada", "La factura electronica no se pudo emitir"},
		{"wallet.low_balance", "Saldo bajo", "El saldo de la billetera esta por debajo del minimo"},
		{"wallet.recharge.failed", "Recarga fallida", "Una recarga de la billetera no se completo"},
	}

	for _, s := range seeds {
		if err := db.Exec(`
INSERT INTO notification_event_types (created_at, updated_at, notification_type_id, event_code, event_name, description, is_active)
VALUES (NOW(), NOW(), ?, ?, ?, ?, true)
ON CONFLICT (notification_type_id, event_code) DO NOTHING
`, typeID, s.code, s.name, s.desc).Error; err != nil {
			return fmt.Errorf("insert assistant event %s: %w", s.code, err)
		}
	}

	if err := db.Exec(`
INSERT INTO notification_event_type_allowed_statuses (notification_event_type_id, order_status_id)
SELECT net.id, os.id
FROM notification_event_types net
JOIN order_statuses os ON os.code IN ('cancelled','failed','rejected','on_hold','refunded','delivery_novelty','delivery_failed','returned','inventory_issue','delivered')
WHERE net.notification_type_id = ? AND net.event_code = 'order.status_changed'
ON CONFLICT DO NOTHING
`, typeID).Error; err != nil {
		return fmt.Errorf("insert assistant allowed statuses: %w", err)
	}

	return nil
}
