package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateProductFamilyParent(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(
		`ALTER TABLE product_families ADD COLUMN IF NOT EXISTS parent_family_id bigint REFERENCES product_families(id)`,
	).Error; err != nil {
		return fmt.Errorf("failed to add product_families.parent_family_id: %w", err)
	}
	if err := r.db.Conn(ctx).Exec(
		`CREATE INDEX IF NOT EXISTS idx_product_families_parent_family_id ON product_families(parent_family_id)`,
	).Error; err != nil {
		return fmt.Errorf("failed to create index on product_families.parent_family_id: %w", err)
	}
	return nil
}
