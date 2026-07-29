package repository

import (
	"context"
	"fmt"
)

const (
	orderErrorsRetentionDays = 30
	orderErrorsDeleteBatch   = 20000
	orderErrorsIncidentMatch = "%cod_includes_shipping%"
)

func (r *Repository) pruneOrderErrors(ctx context.Context) error {
	var deleted int64

	for {
		res := r.db.Conn(ctx).Exec(`
DELETE FROM order_errors
WHERE ctid IN (
    SELECT ctid
    FROM order_errors
    WHERE created_at < NOW() - make_interval(days => ?)
       OR error_message LIKE ?
    LIMIT ?
)`, orderErrorsRetentionDays, orderErrorsIncidentMatch, orderErrorsDeleteBatch)

		if res.Error != nil {
			return fmt.Errorf("prune order_errors: %w", res.Error)
		}

		deleted += res.RowsAffected

		if res.RowsAffected < orderErrorsDeleteBatch {
			break
		}
	}

	if deleted == 0 {
		return nil
	}

	if err := r.db.Conn(ctx).Exec(`VACUUM (FULL, ANALYZE) order_errors`).Error; err != nil {
		fmt.Printf("prune order_errors: se borraron %d filas pero el VACUUM FULL fallo: %v\n", deleted, err)
		return nil
	}

	fmt.Printf("prune order_errors: %d filas eliminadas y tabla compactada\n", deleted)

	return nil
}
