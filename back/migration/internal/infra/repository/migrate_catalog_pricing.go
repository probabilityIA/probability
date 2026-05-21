package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateCatalogPricing(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.ClientGroup{},
		&models.ClientGroupMember{},
		&models.CustomProductPrice{},
		&models.QuantityDiscount{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate catalog pricing: %w", err)
	}

	if err := r.db.Conn(ctx).Exec("DROP TABLE IF EXISTS client_pricing_rule").Error; err != nil {
		return fmt.Errorf("failed to drop client_pricing_rule: %w", err)
	}

	return nil
}
