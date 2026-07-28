package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateProductChannelCategories(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS channel_categories jsonb`,
	).Error; err != nil {
		return fmt.Errorf("failed to add products.channel_categories: %w", err)
	}
	return nil
}
