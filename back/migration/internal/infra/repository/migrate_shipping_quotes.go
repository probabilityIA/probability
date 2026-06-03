package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateShippingQuotes(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.ShippingQuote{}); err != nil {
		return fmt.Errorf("automigrate shipping_quotes: %w", err)
	}
	return nil
}
