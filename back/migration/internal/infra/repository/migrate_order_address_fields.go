package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateOrderAddressFields(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.Order{}); err != nil {
		return fmt.Errorf("failed to auto-migrate order address fields: %w", err)
	}

	return nil
}
