package repository

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/add_shipments_geozone.sql
var addShipmentsGeozoneSQL string

func (r *Repository) migrateShipmentsGeozone(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(addShipmentsGeozoneSQL).Error; err != nil {
		return fmt.Errorf("failed to migrate shipments geozone: %w", err)
	}
	return nil
}
