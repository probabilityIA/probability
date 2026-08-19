package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateBusinessModuleOverrideExpiry(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.BusinessModuleOverride{}); err != nil {
		return fmt.Errorf("failed to auto-migrate business_module_overrides expiry: %w", err)
	}

	return nil
}
