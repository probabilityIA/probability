package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateNavbarContent(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`ALTER TABLE business_website_config ADD COLUMN IF NOT EXISTS navbar_content jsonb`).Error; err != nil {
		return fmt.Errorf("add navbar_content to business_website_config: %w", err)
	}

	return nil
}
