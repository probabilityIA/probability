package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateOrderStatusChangedAt(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE orders
ADD COLUMN IF NOT EXISTS status_changed_at TIMESTAMPTZ
`).Error; err != nil {
		return fmt.Errorf("add status_changed_at to orders: %w", err)
	}

	if err := db.Exec(`
UPDATE orders o
SET status_changed_at = h.last_change
FROM (
    SELECT order_id, MAX(created_at) AS last_change
    FROM order_history
    WHERE deleted_at IS NULL
    GROUP BY order_id
) h
WHERE h.order_id = o.id
  AND o.status_changed_at IS NULL
  AND o.status_changed_by <> ''
`).Error; err != nil {
		return fmt.Errorf("backfill status_changed_at from order_history: %w", err)
	}

	return nil
}
