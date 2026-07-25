package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateAccountingInvoices(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(
		&models.AccountingInvoice{},
		&models.AccountingInvoiceItem{},
	); err != nil {
		return fmt.Errorf("automigrate accounting invoices: %w", err)
	}

	if err := db.Exec(`
UPDATE accounting_invoices SET dian_status = 'NONE' WHERE dian_status IS NULL OR dian_status = ''
`).Error; err != nil {
		return fmt.Errorf("backfill dian_status: %w", err)
	}

	if err := db.Exec(`
INSERT INTO accounting_concepts (code, name, description, kind, is_real_income, is_automatic, source_type, is_active, created_at, updated_at)
VALUES ('SERVICE_INVOICE', 'Facturacion de servicios', 'Facturas emitidas por la plataforma a negocios', 'INCOME', TRUE, FALSE, 'MANUAL', TRUE, NOW(), NOW())
ON CONFLICT (code) DO NOTHING
`).Error; err != nil {
		return fmt.Errorf("seed concept SERVICE_INVOICE: %w", err)
	}

	return nil
}
