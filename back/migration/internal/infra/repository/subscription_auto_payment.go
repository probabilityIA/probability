package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateSubscriptionAutoPayment(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.Business{}); err != nil {
		return fmt.Errorf("failed to auto-migrate subscription auto payment: %w", err)
	}

	return nil
}
