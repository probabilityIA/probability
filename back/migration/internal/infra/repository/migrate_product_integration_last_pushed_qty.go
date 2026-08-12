package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateProductIntegrationLastPushedQty(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.ProductBusinessIntegration{}); err != nil {
		return fmt.Errorf("failed to auto-migrate product_business_integrations: %w", err)
	}
	return nil
}
