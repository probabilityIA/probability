package repository

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/backfill_orders_geozone_by_point.sql
var backfillOrdersGeozoneByPointSQL string

func (r *Repository) backfillOrdersGeozoneByPoint(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(backfillOrdersGeozoneByPointSQL).Error; err != nil {
		return fmt.Errorf("failed to backfill orders geozone by point: %w", err)
	}
	return nil
}
