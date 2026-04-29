package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateWarehouseStructureType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE warehouses
    ADD COLUMN IF NOT EXISTS structure_type VARCHAR(20) NOT NULL DEFAULT 'simple'
`).Error; err != nil {
		return fmt.Errorf("add structure_type column: %w", err)
	}

	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_warehouses_structure_type ON warehouses(structure_type)`).Error; err != nil {
		return fmt.Errorf("create structure_type index: %w", err)
	}

	if err := db.Exec(`
UPDATE warehouses w
SET structure_type = 'wms'
WHERE w.deleted_at IS NULL
  AND w.structure_type = 'simple'
  AND EXISTS (
      SELECT 1
      FROM warehouse_racks r
      JOIN warehouse_aisles a ON a.id = r.aisle_id AND a.deleted_at IS NULL
      JOIN warehouse_zones z ON z.id = a.zone_id AND z.deleted_at IS NULL
      WHERE z.warehouse_id = w.id AND r.deleted_at IS NULL
  )
`).Error; err != nil {
		return fmt.Errorf("backfill wms warehouses: %w", err)
	}

	if err := db.Exec(`
UPDATE warehouses w
SET structure_type = 'zones'
WHERE w.deleted_at IS NULL
  AND w.structure_type = 'simple'
  AND EXISTS (
      SELECT 1
      FROM warehouse_zones z
      WHERE z.warehouse_id = w.id AND z.deleted_at IS NULL
  )
`).Error; err != nil {
		return fmt.Errorf("backfill zones warehouses: %w", err)
	}

	return nil
}
