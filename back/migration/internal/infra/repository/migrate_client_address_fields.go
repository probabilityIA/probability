package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateClientAddressFields(ctx context.Context) error {
	db := r.db.Conn(ctx)
	if err := db.AutoMigrate(&models.Client{}); err != nil {
		return fmt.Errorf("failed to auto-migrate clients (address/city/notes): %w", err)
	}
	return nil
}
