package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateShipmentSyncLogs(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.ShipmentSyncLog{}); err != nil {
		return fmt.Errorf("failed to auto-migrate shipment_sync_logs table: %w", err)
	}
	return nil
}
