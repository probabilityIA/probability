package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateRackSide(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.WarehouseRack{}); err != nil {
		return fmt.Errorf("failed to auto-migrate warehouse_racks side column: %w", err)
	}
	return nil
}
