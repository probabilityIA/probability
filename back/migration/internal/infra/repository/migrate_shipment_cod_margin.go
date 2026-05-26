package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateShipmentCodMargin(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Shipment{}); err != nil {
		return fmt.Errorf("extend shipments with cod_customer_charge, cod_applied_margin: %w", err)
	}
	return nil
}
