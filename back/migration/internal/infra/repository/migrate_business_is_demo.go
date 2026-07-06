package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateBusinessIsDemo(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Business{}); err != nil {
		return fmt.Errorf("failed to auto-migrate business is_demo: %w", err)
	}
	return nil
}
