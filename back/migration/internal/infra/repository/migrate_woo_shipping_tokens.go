package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWooShippingTokens(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.WooShippingToken{}); err != nil {
		return fmt.Errorf("failed to auto-migrate woo_shipping_tokens: %w", err)
	}
	return nil
}
