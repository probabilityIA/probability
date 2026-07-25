package repository

import (
	"context"
	"fmt"
)

const wooCommerceIntegrationTypeID = 4

func (r *Repository) migrateWooCommerceStatusMappings(ctx context.Context) error {
	mappings := []struct {
		original string
		target   string
		desc     string
	}{
		{"pending", "pending", "Pedido creado, esperando pago"},
		{"checkout-draft", "pending", "Borrador de checkout en WooCommerce"},
		{"processing", "processing", "Pago recibido, pedido en preparacion"},
		{"addi-approved", "processing", "Financiacion aprobada por Addi"},
		{"on-hold", "on_hold", "Pedido en espera de confirmacion de pago"},
		{"completed", "completed", "Pedido despachado y completado"},
		{"cancelled", "cancelled", "Pedido cancelado"},
		{"trash", "cancelled", "Pedido enviado a la papelera"},
		{"refunded", "refunded", "Pedido reembolsado"},
		{"failed", "failed", "Pago fallido o rechazado"},
	}

	const upsert = `
		INSERT INTO order_status_mappings
			(created_at, updated_at, integration_type_id, original_status, order_status_id, is_active, description)
		SELECT NOW(), NOW(), ?, ?, os.id, true, ?
		FROM order_statuses os
		WHERE os.code = ? AND os.deleted_at IS NULL
		ON CONFLICT (integration_type_id, original_status) DO NOTHING`

	for _, m := range mappings {
		err := r.db.Conn(ctx).Exec(upsert, wooCommerceIntegrationTypeID, m.original, m.desc, m.target).Error
		if err != nil {
			return fmt.Errorf("seed mapeo woocommerce %s: %w", m.original, err)
		}
	}

	return nil
}
