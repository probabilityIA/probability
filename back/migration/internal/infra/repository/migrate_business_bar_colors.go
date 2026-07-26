package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateBusinessBarColors(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE business
ADD COLUMN IF NOT EXISTS sidebar_color VARCHAR(7)
`).Error; err != nil {
		return fmt.Errorf("add sidebar_color to business: %w", err)
	}

	if err := db.Exec(`
ALTER TABLE business
ADD COLUMN IF NOT EXISTS topbar_color VARCHAR(7)
`).Error; err != nil {
		return fmt.Errorf("add topbar_color to business: %w", err)
	}

	return nil
}
