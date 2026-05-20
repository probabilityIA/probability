package repository

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/backfill_orders_geozone.sql
var backfillOrdersGeozoneSQL string

func (r *Repository) backfillOrdersGeozone(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(backfillOrdersGeozoneSQL).Error; err != nil {
		return fmt.Errorf("failed to backfill orders geozone: %w", err)
	}
	return nil
}
