package repository

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/add_integrations_store_id_unique.sql
var addIntegrationsStoreIDUniqueSQL string

func (r *Repository) migrateIntegrationsStoreIDUnique(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(addIntegrationsStoreIDUniqueSQL).Error; err != nil {
		return fmt.Errorf("failed to migrate integrations store_id unique index: %w", err)
	}
	return nil
}
