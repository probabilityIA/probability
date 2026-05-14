package repository

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/create_geozones.sql
var createGeozonesSQL string

func (r *Repository) migrateGeozones(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(createGeozonesSQL).Error; err != nil {
		return fmt.Errorf("failed to migrate geozones: %w", err)
	}
	return nil
}
