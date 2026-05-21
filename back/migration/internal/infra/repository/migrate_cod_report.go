package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateCodReport(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.CarrierCodConfig{},
		&models.CodPaymentCut{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate cod report: %w", err)
	}
	return nil
}
