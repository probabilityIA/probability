package repository

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/add_geozone_analytics.sql
var addGeozoneAnalyticsSQL string

func (r *Repository) migrateGeozoneAnalytics(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(addGeozoneAnalyticsSQL).Error; err != nil {
		return fmt.Errorf("failed to migrate geozone analytics: %w", err)
	}
	return nil
}
