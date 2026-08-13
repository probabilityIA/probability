package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateSiigoReferrals(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.SiigoReferral{}); err != nil {
		return fmt.Errorf("automigrate siigo_referrals: %w", err)
	}
	return nil
}
