package repository

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/create_inventory_submodule_resources.sql
var inventorySubmoduleResourcesSQL string

func (r *Repository) migrateInventorySubmoduleResources(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(inventorySubmoduleResourcesSQL).Error; err != nil {
		return fmt.Errorf("failed to seed inventory submodule resources: %w", err)
	}
	return nil
}
