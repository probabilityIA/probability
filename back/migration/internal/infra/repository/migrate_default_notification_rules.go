package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateDefaultNotificationRules(ctx context.Context) error {
	db := r.db.Conn(ctx)

	var eventTypeID uint
	if err := db.Raw(`
SELECT net.id FROM notification_event_types net
JOIN notification_types nt ON nt.id = net.notification_type_id
WHERE nt.code = 'assistant' AND net.event_code = 'order.status_changed' AND net.deleted_at IS NULL
LIMIT 1`).Scan(&eventTypeID).Error; err != nil {
		return fmt.Errorf("buscar evento de cambio de estado del asistente: %w", err)
	}
	if eventTypeID == 0 {
		return fmt.Errorf("no existe el evento order.status_changed del canal assistant: correr migrateAssistantAlerts primero")
	}

	if err := db.Exec(`
INSERT INTO business_notification_configs (created_at, updated_at, business_id, integration_id, notification_type_id, notification_event_type_id, enabled, description)
SELECT NOW(), NOW(), b.id, plataforma.id, ?, ?, true, 'Regla predeterminada del sistema'
FROM business b
JOIN LATERAL (
	SELECT i.id FROM integrations i
	JOIN integration_types it ON it.id = i.integration_type_id
	WHERE i.business_id = b.id AND it.code = 'platform' AND i.deleted_at IS NULL
	ORDER BY i.id LIMIT 1
) plataforma ON true
WHERE b.deleted_at IS NULL
AND NOT EXISTS (
	SELECT 1 FROM business_notification_configs c
	WHERE c.business_id = b.id AND c.notification_type_id = ? AND c.notification_event_type_id = ? AND c.deleted_at IS NULL
)`, assistantNotificationTypeID, eventTypeID, assistantNotificationTypeID, eventTypeID).Error; err != nil {
		return fmt.Errorf("crear reglas predeterminadas del asistente: %w", err)
	}

	if err := db.Exec(`
UPDATE business_notification_configs SET enabled = true, updated_at = NOW()
WHERE notification_type_id = ? AND notification_event_type_id = ? AND deleted_at IS NULL`, assistantNotificationTypeID, eventTypeID).Error; err != nil {
		return fmt.Errorf("activar reglas del asistente: %w", err)
	}

	if err := db.Exec(`
DELETE FROM business_notification_config_order_statuses
WHERE business_notification_config_id IN (
	SELECT id FROM business_notification_configs
	WHERE notification_type_id = ? AND notification_event_type_id = ? AND deleted_at IS NULL
)`, assistantNotificationTypeID, eventTypeID).Error; err != nil {
		return fmt.Errorf("limpiar estados de las reglas del asistente: %w", err)
	}

	if err := db.Exec(`
INSERT INTO business_notification_config_order_statuses (business_notification_config_id, order_status_id, created_at)
SELECT c.id, os.id, NOW()
FROM business_notification_configs c
JOIN order_statuses os ON os.code IN ('cancelled', 'cancel_requested')
WHERE c.notification_type_id = ? AND c.notification_event_type_id = ? AND c.deleted_at IS NULL
ON CONFLICT DO NOTHING`, assistantNotificationTypeID, eventTypeID).Error; err != nil {
		return fmt.Errorf("asignar estados a las reglas del asistente: %w", err)
	}

	return nil
}
