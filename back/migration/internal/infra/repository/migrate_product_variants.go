package repository

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateProductVariants(ctx context.Context) error {
	return r.db.Conn(ctx).AutoMigrate(
		&models.ProductFamily{},
		&models.Product{},
		&models.ProductBusinessIntegration{},
	)
}
