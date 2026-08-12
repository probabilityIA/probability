package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateProductIntegrationVariantUnique(ctx context.Context) error {
	conn := r.db.Conn(ctx)

	if err := conn.Exec(`DROP INDEX IF EXISTS idx_product_integration`).Error; err != nil {
		return fmt.Errorf("dropping idx_product_integration: %w", err)
	}

	stmt := `CREATE UNIQUE INDEX IF NOT EXISTS idx_product_integration
		ON product_business_integrations (product_id, integration_id, COALESCE(external_variant_id, ''))`
	if err := conn.Exec(stmt).Error; err != nil {
		return fmt.Errorf("creating idx_product_integration: %w", err)
	}

	return nil
}
