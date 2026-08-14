package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateInventoryCompareSnapshot(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.InventoryCompareSnapshot{}); err != nil {
		return fmt.Errorf("failed to auto-migrate inventory_compare_snapshots: %w", err)
	}

	if err := db.Exec(`DROP INDEX IF EXISTS idx_inv_compare_unique`).Error; err != nil {
		return fmt.Errorf("failed to drop idx_inv_compare_unique: %w", err)
	}

	if err := db.Exec(`ALTER TABLE inventory_compare_snapshots ALTER COLUMN external_variant_id SET DEFAULT ''`).Error; err != nil {
		return fmt.Errorf("failed to default external_variant_id: %w", err)
	}
	if err := db.Exec(`UPDATE inventory_compare_snapshots SET external_variant_id = '' WHERE external_variant_id IS NULL`).Error; err != nil {
		return fmt.Errorf("failed to backfill external_variant_id: %w", err)
	}
	if err := db.Exec(`ALTER TABLE inventory_compare_snapshots ALTER COLUMN external_variant_id SET NOT NULL`).Error; err != nil {
		return fmt.Errorf("failed to set external_variant_id not null: %w", err)
	}

	sql := `CREATE UNIQUE INDEX IF NOT EXISTS idx_inv_compare_unique
		ON inventory_compare_snapshots (integration_id, product_id, external_variant_id)`
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to create idx_inv_compare_unique: %w", err)
	}

	return nil
}
