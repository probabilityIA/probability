package repository

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/add_shipping_neighborhood.sql
var addShippingNeighborhoodSQL string

func (r *Repository) migrateShippingNeighborhood(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(addShippingNeighborhoodSQL).Error; err != nil {
		return fmt.Errorf("failed to migrate shipping_neighborhood: %w", err)
	}
	return nil
}
