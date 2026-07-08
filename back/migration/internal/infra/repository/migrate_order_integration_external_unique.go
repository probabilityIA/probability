package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateOrderIntegrationExternalUnique(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`DROP INDEX IF EXISTS idx_integration_external_id`).Error; err != nil {
		return fmt.Errorf("drop idx_integration_external_id: %w", err)
	}

	if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_integration_external_id ON orders (integration_id, external_id)`).Error; err != nil {
		return fmt.Errorf("create composite idx_integration_external_id: %w", err)
	}

	return nil
}
