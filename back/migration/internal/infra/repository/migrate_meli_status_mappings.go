package repository

import (
	"context"
	"fmt"
)

const meliIntegrationTypeID = 3

func (r *Repository) migrateMeliStatusMappings(ctx context.Context) error {
	mappings := []struct {
		original string
		target   string
		desc     string
	}{
		{"pending", "pending", "Orden creada, esperando pago"},
		{"confirmed", "pending", "Orden confirmada en MercadoLibre, aun sin pago"},
		{"payment_required", "pending", "MercadoLibre espera el pago del comprador"},
		{"payment_in_process", "pending", "Pago en proceso de aprobacion"},
		{"paid", "processing", "Pago aprobado, entra al flujo de almacen"},
		{"partially_paid", "processing", "Pago parcial recibido, entra al flujo"},
		{"cancelled", "cancelled", "Orden cancelada"},
		{"invalid", "cancelled", "Orden invalidada por MercadoLibre"},
	}

	const upsert = `
		INSERT INTO order_status_mappings
			(created_at, updated_at, integration_type_id, original_status, order_status_id, is_active, description)
		SELECT NOW(), NOW(), ?, ?, os.id, true, ?
		FROM order_statuses os
		WHERE os.code = ? AND os.deleted_at IS NULL
		ON CONFLICT (integration_type_id, original_status) DO NOTHING`

	for _, m := range mappings {
		err := r.db.Conn(ctx).Exec(upsert, meliIntegrationTypeID, m.original, m.desc, m.target).Error
		if err != nil {
			return fmt.Errorf("seed mapeo mercadolibre %s: %w", m.original, err)
		}
	}

	const backfill = `
		UPDATE orders o
		SET status_id = m.order_status_id,
		    updated_at = NOW()
		FROM order_status_mappings m
		WHERE m.integration_type_id = ?
		  AND m.deleted_at IS NULL
		  AND m.is_active = true
		  AND o.platform = 'mercadolibre'
		  AND o.status = m.original_status
		  AND o.status_id IS NULL`

	if err := r.db.Conn(ctx).Exec(backfill, meliIntegrationTypeID).Error; err != nil {
		return fmt.Errorf("backfill status_id de ordenes mercadolibre: %w", err)
	}

	return nil
}
