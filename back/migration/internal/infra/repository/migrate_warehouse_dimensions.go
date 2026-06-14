package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWarehouseDimensions(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(
		&models.WarehouseAisle{},
		&models.WarehouseRack{},
		&models.WarehouseLayout{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate warehouse dimensions: %w", err)
	}

	return nil
}
