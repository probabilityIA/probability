package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateInventoryOperations(ctx context.Context) error {
	db := r.db.Conn(ctx)
	if err := db.AutoMigrate(
		&models.PutawayRule{},
		&models.PutawaySuggestion{},
		&models.ReplenishmentTask{},
		&models.CrossDockLink{},
		&models.ProductVelocity{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate inventory operations tables: %w", err)
	}
	return nil
}
