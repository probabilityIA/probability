package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateMarketingLeads(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.MarketingLead{}); err != nil {
		return fmt.Errorf("automigrate marketing_leads: %w", err)
	}
	return nil
}
