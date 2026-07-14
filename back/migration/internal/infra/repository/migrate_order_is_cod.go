package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateOrderIsCod(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Order{}); err != nil {
		return fmt.Errorf("add orders.is_cod: %w", err)
	}

	backfill := `UPDATE orders SET is_cod = true WHERE is_cod = false AND cod_total IS NOT NULL AND cod_total > 0`
	if err := r.db.Conn(ctx).Exec(backfill).Error; err != nil {
		return fmt.Errorf("backfill orders.is_cod: %w", err)
	}

	deactivateCod := `UPDATE payment_methods SET is_active = false WHERE code = 'cod' AND is_active = true`
	if err := r.db.Conn(ctx).Exec(deactivateCod).Error; err != nil {
		return fmt.Errorf("deactivate cod payment method: %w", err)
	}

	return nil
}
