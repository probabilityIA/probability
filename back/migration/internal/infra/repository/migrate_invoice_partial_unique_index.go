package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateInvoicePartialUniqueIndex(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`DROP INDEX IF EXISTS idx_order_provider`).Error; err != nil {
		return fmt.Errorf("drop idx_order_provider: %w", err)
	}

	if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_order_provider ON invoices (order_id, invoicing_provider_id) WHERE status <> 'cancelled' AND deleted_at IS NULL`).Error; err != nil {
		return fmt.Errorf("create partial idx_order_provider: %w", err)
	}

	return nil
}
