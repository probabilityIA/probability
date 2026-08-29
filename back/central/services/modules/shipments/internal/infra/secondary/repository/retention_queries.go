package repository

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) PurgeSyncLogsOlderThan(ctx context.Context, days, batchSize int) (int64, error) {
	if days <= 0 {
		return 0, fmt.Errorf("invalid retention days: %d", days)
	}
	if batchSize <= 0 {
		batchSize = 1000
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	result := r.db.Conn(ctx).Exec(`
		DELETE FROM shipment_sync_logs
		WHERE id IN (
			SELECT id FROM shipment_sync_logs
			WHERE created_at < ?
			ORDER BY id
			LIMIT ?
		)`, cutoff, batchSize)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
