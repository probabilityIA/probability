package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateInventoryCapture(ctx context.Context) error {
	db := r.db.Conn(ctx)
	if err := db.AutoMigrate(
		&models.LicensePlate{},
		&models.LicensePlateLine{},
		&models.ScanEvent{},
		&models.InventorySyncLog{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate inventory capture tables: %w", err)
	}
	return nil
}
