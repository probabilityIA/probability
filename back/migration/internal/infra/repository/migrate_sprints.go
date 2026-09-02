package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateSprints(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Sprint{}); err != nil {
		return fmt.Errorf("failed to auto-migrate sprints: %w", err)
	}
	if err := r.db.Conn(ctx).AutoMigrate(&models.Ticket{}); err != nil {
		return fmt.Errorf("failed to add sprint_id to tickets: %w", err)
	}
	return nil
}
