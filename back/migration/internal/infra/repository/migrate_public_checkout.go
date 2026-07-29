package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migratePublicCheckout(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.PublicCheckout{}); err != nil {
		return fmt.Errorf("automigrate public_checkouts: %w", err)
	}
	return nil
}
