package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWalletTxIntegration(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.WalletTransaction{}); err != nil {
		return fmt.Errorf("failed to auto-migrate wallet transaction integration columns: %w", err)
	}
	return nil
}
