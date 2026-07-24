package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateFiscalProfiles(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(
		&models.BusinessFiscalProfile{},
		&models.AccountingTax{},
		&models.AccountingInvoice{},
	); err != nil {
		return fmt.Errorf("automigrate fiscal profiles: %w", err)
	}

	kinds := []struct{ code, kind string }{
		{"IVA", "CHARGE"},
		{"RETEFUENTE", "WITHHOLDING"},
		{"RETEICA", "WITHHOLDING"},
		{"GMF_4X1000", "OTHER"},
	}
	for _, k := range kinds {
		if err := db.Exec(`UPDATE accounting_taxes SET kind = ? WHERE code = ?`, k.kind, k.code).Error; err != nil {
			return fmt.Errorf("set tax kind %s: %w", k.code, err)
		}
	}

	if err := db.Exec(`UPDATE accounting_taxes SET kind = 'CHARGE' WHERE kind IS NULL OR kind = ''`).Error; err != nil {
		return fmt.Errorf("backfill tax kind: %w", err)
	}

	return nil
}
