package repository

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

//go:embed sql/integration_stats_triggers.sql
var integrationStatsTriggersSQL string

func (r *Repository) migrateIntegrationStats(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.IntegrationStat{}); err != nil {
		return fmt.Errorf("failed to automigrate integration_stats: %w", err)
	}
	if err := r.db.Conn(ctx).Exec(integrationStatsTriggersSQL).Error; err != nil {
		return fmt.Errorf("failed to create integration_stats triggers: %w", err)
	}
	return nil
}
