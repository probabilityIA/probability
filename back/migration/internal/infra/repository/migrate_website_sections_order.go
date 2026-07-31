package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateWebsiteSectionsOrder(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE business_website_config
    ADD COLUMN IF NOT EXISTS sections_order JSONB
`).Error; err != nil {
		return fmt.Errorf("add sections_order to business_website_config: %w", err)
	}

	if err := db.Exec(`
ALTER TABLE business_website_config
    ADD COLUMN IF NOT EXISTS theme_content JSONB
`).Error; err != nil {
		return fmt.Errorf("add theme_content to business_website_config: %w", err)
	}

	if err := db.Exec(`
ALTER TABLE business_website_config
    ADD COLUMN IF NOT EXISTS hidden_categories JSONB
`).Error; err != nil {
		return fmt.Errorf("add hidden_categories to business_website_config: %w", err)
	}

	return nil
}
