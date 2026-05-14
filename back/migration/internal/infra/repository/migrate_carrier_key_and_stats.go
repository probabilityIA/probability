package repository

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/add_carrier_key_and_stats.sql
var addCarrierKeyAndStatsSQL string

func (r *Repository) migrateCarrierKeyAndStats(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(addCarrierKeyAndStatsSQL).Error; err != nil {
		return fmt.Errorf("failed to migrate carrier_key + geozone_carrier_stats: %w", err)
	}
	return nil
}
