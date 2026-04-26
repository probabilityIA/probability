package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateShipmentDestinationGeo(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Shipment{}); err != nil {
		return fmt.Errorf("failed to auto-migrate shipment for destination geo: %w", err)
	}
	return nil
}
