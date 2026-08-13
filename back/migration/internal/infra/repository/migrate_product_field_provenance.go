package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateProductFieldProvenance(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.ProductFieldOrigin{}); err != nil {
		return fmt.Errorf("failed to auto-migrate product_field_origins: %w", err)
	}
	if err := r.db.Conn(ctx).AutoMigrate(&models.ProductFieldChange{}); err != nil {
		return fmt.Errorf("failed to auto-migrate product_field_changes: %w", err)
	}
	if err := r.db.Conn(ctx).AutoMigrate(&models.ProductBusinessIntegration{}); err != nil {
		return fmt.Errorf("failed to auto-migrate product_business_integrations: %w", err)
	}
	return nil
}
