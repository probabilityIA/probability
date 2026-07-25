package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateAccountingServices(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(
		&models.AccountingService{},
		&models.AccountingInvoiceItem{},
	); err != nil {
		return fmt.Errorf("automigrate accounting services: %w", err)
	}

	if err := db.Exec(`
INSERT INTO accounting_services (code, name, unspsc_code, unit_price, is_active, created_at, updated_at)
VALUES ('SOFT-LOG', 'Alquiler de software logistico', '', 0, TRUE, NOW(), NOW())
ON CONFLICT (code) DO NOTHING
`).Error; err != nil {
		return fmt.Errorf("seed accounting service: %w", err)
	}

	return nil
}
