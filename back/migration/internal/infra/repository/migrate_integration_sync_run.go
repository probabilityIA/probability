package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateIntegrationSyncRuns(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.IntegrationSyncRun{}); err != nil {
		return fmt.Errorf("failed to auto-migrate integration_sync_runs: %w", err)
	}
	return nil
}
