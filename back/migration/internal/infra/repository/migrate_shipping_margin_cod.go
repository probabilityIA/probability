package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateShippingMarginCOD(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.ShippingMargin{}); err != nil {
		return fmt.Errorf("extend shipping_margin with cod margin columns: %w", err)
	}

	if err := r.db.Conn(ctx).Exec(`
ALTER TABLE shipping_margin
ADD COLUMN IF NOT EXISTS cod_margin_amount decimal(15,2) NOT NULL DEFAULT 0
`).Error; err != nil {
		return fmt.Errorf("add shipping_margin.cod_margin_amount: %w", err)
	}

	return nil
}
