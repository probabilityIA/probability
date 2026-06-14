package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWarehouseLayout(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.WarehouseLayout{}); err != nil {
		return fmt.Errorf("failed to auto-migrate warehouse_layouts table: %w", err)
	}

	return nil
}
