package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateSubscriptionCutoffDay(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.Business{}); err != nil {
		return fmt.Errorf("failed to auto-migrate subscription cutoff day: %w", err)
	}

	return nil
}
