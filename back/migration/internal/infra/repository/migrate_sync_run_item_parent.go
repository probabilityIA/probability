package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateSyncRunItemParent(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.IntegrationSyncRunItem{}); err != nil {
		return fmt.Errorf("failed to auto-migrate integration_sync_run_items: %w", err)
	}
	return nil
}
