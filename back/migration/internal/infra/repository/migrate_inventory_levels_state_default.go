package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateInventoryLevelsStateDefault(ctx context.Context) error {
	db := r.db.Conn(ctx)

	var availableID uint
	if err := db.Raw(`SELECT id FROM inventory_states WHERE code = 'available' LIMIT 1`).Scan(&availableID).Error; err != nil {
		return fmt.Errorf("lookup available state: %w", err)
	}
	if availableID == 0 {
		return fmt.Errorf("inventory_states 'available' not seeded")
	}

	mergeSQL := `
WITH null_rows AS (
    SELECT id, product_id, warehouse_id, business_id, location_id, lot_id, quantity, reserved_qty
    FROM inventory_levels
    WHERE deleted_at IS NULL AND state_id IS NULL
),
sibling AS (
    SELECT n.id AS null_id, l.id AS available_id, n.quantity AS add_qty, n.reserved_qty AS add_reserved
    FROM null_rows n
    JOIN inventory_levels l
      ON l.deleted_at IS NULL
     AND l.state_id = ?
     AND l.product_id = n.product_id
     AND l.warehouse_id = n.warehouse_id
     AND l.business_id = n.business_id
     AND l.location_id IS NOT DISTINCT FROM n.location_id
     AND l.lot_id IS NOT DISTINCT FROM n.lot_id
),
merged AS (
    UPDATE inventory_levels l
    SET quantity = l.quantity + s.add_qty,
        reserved_qty = l.reserved_qty + s.add_reserved,
        available_qty = (l.quantity + s.add_qty) - (l.reserved_qty + s.add_reserved),
        updated_at = NOW()
    FROM sibling s
    WHERE l.id = s.available_id
    RETURNING s.null_id
)
DELETE FROM inventory_levels WHERE id IN (SELECT null_id FROM merged);
`
	if err := db.Exec(mergeSQL, availableID).Error; err != nil {
		return fmt.Errorf("merge null state rows: %w", err)
	}

	if err := db.Exec(`
UPDATE inventory_levels
SET state_id = ?, available_qty = quantity - reserved_qty, updated_at = NOW()
WHERE deleted_at IS NULL AND state_id IS NULL
`, availableID).Error; err != nil {
		return fmt.Errorf("backfill remaining null state rows: %w", err)
	}

	if err := db.Exec(fmt.Sprintf(`ALTER TABLE inventory_levels ALTER COLUMN state_id SET DEFAULT %d`, availableID)).Error; err != nil {
		return fmt.Errorf("set state_id default: %w", err)
	}

	if err := db.Exec(`ALTER TABLE inventory_levels ALTER COLUMN state_id SET NOT NULL`).Error; err != nil {
		return fmt.Errorf("set state_id not null: %w", err)
	}

	return nil
}
