package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateShipmentProbabilityGuide(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Shipment{}); err != nil {
		return fmt.Errorf("add probability_guide_url to shipments: %w", err)
	}
	return nil
}
