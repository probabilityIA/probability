package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateOrderGeoConfidence(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Order{}); err != nil {
		return fmt.Errorf("failed to auto-migrate order for geo confidence: %w", err)
	}
	return nil
}
