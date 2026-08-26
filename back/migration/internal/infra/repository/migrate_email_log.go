package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateEmailLogGeneric(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.EmailLog{}); err != nil {
		return fmt.Errorf("failed to auto-migrate email_logs: %w", err)
	}
	return nil
}
