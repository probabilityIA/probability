package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateClientMultiBusiness(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`DROP INDEX IF EXISTS idx_clients_user_id`).Error; err != nil {
		return fmt.Errorf("drop idx_clients_user_id: %w", err)
	}

	if err := db.Exec(`
CREATE UNIQUE INDEX IF NOT EXISTS idx_business_client_user
    ON client (business_id, user_id)
    WHERE user_id IS NOT NULL AND deleted_at IS NULL
`).Error; err != nil {
		return fmt.Errorf("create idx_business_client_user: %w", err)
	}

	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_client_user_id ON client (user_id)`).Error; err != nil {
		return fmt.Errorf("create idx_client_user_id: %w", err)
	}

	if err := db.Exec(`ALTER TABLE customer_address ADD COLUMN IF NOT EXISTS is_primary BOOLEAN NOT NULL DEFAULT false`).Error; err != nil {
		return fmt.Errorf("add is_primary to customer_address: %w", err)
	}

	return nil
}
