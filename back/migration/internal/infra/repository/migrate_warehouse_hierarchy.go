package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWarehouseHierarchy(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(
		&models.WarehouseZone{},
		&models.WarehouseAisle{},
		&models.WarehouseRack{},
		&models.WarehouseRackLevel{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate warehouse hierarchy tables: %w", err)
	}

	if err := db.AutoMigrate(&models.WarehouseLocation{}); err != nil {
		return fmt.Errorf("failed to auto-migrate warehouse_locations with hierarchy columns: %w", err)
	}

	return nil
}
