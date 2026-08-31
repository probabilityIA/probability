package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateUserTourProgress(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.UserTourProgress{}); err != nil {
		return fmt.Errorf("failed to auto-migrate user_tour_progress: %w", err)
	}
	return nil
}
