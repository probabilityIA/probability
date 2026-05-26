package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateShippingMarginCOD(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.ShippingMargin{}); err != nil {
		return fmt.Errorf("extend shipping_margin with cod_margin_percent: %w", err)
	}
	return nil
}
