package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateOrderStatusSource(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Order{}, &models.OrderHistory{}); err != nil {
		return fmt.Errorf("add status_source columns: %w", err)
	}

	backfillHistory := `UPDATE order_history SET source = 'user' WHERE source IS NULL OR source = ''`
	if err := r.db.Conn(ctx).Exec(backfillHistory).Error; err != nil {
		return fmt.Errorf("backfill order_history.source: %w", err)
	}

	return nil
}
