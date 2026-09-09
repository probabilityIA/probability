package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

type pushEventSeed struct {
	code string
	name string
	desc string
}

func (r *Repository) migratePushNotifications(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.DeviceToken{}); err != nil {
		return fmt.Errorf("automigrate device_tokens: %w", err)
	}

	if err := db.Exec(`
INSERT INTO notification_types (created_at, updated_at, name, code, description, icon, is_active)
VALUES (NOW(), NOW(), 'Push', 'push', 'Notificaciones push a la aplicacion movil', 'bell', true)
ON CONFLICT (code) DO NOTHING
`).Error; err != nil {
		return fmt.Errorf("insert push notification type: %w", err)
	}

	var pushTypeID uint
	if err := db.Raw(`SELECT id FROM notification_types WHERE code = 'push' LIMIT 1`).Scan(&pushTypeID).Error; err != nil {
		return fmt.Errorf("lookup push notification type: %w", err)
	}
	if pushTypeID == 0 {
		return fmt.Errorf("push notification type no quedo creado")
	}

	seeds := []pushEventSeed{
		{"shipment.guide_generated", "Guia de envio generada", "Avisa al negocio cuando la transportadora devuelve la guia"},
		{"order.status_changed", "Cambio de estado de la orden", "Cubre novedad de entrega, entregado y devuelto"},
		{"order.delivered", "Pedido entregado", "Cierra el ciclo del envio"},
		{"wallet.low_balance", "Saldo bajo", "Avisa antes de que el saldo bloquee la generacion de guias"},
	}

	for _, s := range seeds {
		if err := db.Exec(`
INSERT INTO notification_event_types (created_at, updated_at, notification_type_id, event_code, event_name, description, is_active)
VALUES (NOW(), NOW(), ?, ?, ?, ?, true)
ON CONFLICT (notification_type_id, event_code) DO NOTHING
`, pushTypeID, s.code, s.name, s.desc).Error; err != nil {
			return fmt.Errorf("insert push event %s: %w", s.code, err)
		}
	}

	return nil
}
