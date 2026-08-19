package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateSubscriptionTypeOverage(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.SubscriptionType{}); err != nil {
		return fmt.Errorf("failed to auto-migrate subscription_types overage fields: %w", err)
	}

	return nil
}
