package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateCommercialProspects(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.CommercialProspect{}); err != nil {
		return fmt.Errorf("automigrate commercial_prospects: %w", err)
	}
	return nil
}
