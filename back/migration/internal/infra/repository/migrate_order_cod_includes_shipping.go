package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateOrderCodIncludesShipping(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Order{}); err != nil {
		return fmt.Errorf("add orders.cod_includes_shipping: %w", err)
	}

	backfill := `UPDATE orders SET cod_includes_shipping = false WHERE platform = 'manual' AND cod_includes_shipping = true`
	if err := r.db.Conn(ctx).Exec(backfill).Error; err != nil {
		return fmt.Errorf("backfill orders.cod_includes_shipping: %w", err)
	}

	return nil
}
