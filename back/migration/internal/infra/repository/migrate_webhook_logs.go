package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWebhookLogs(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.WebhookLog{}); err != nil {
		return fmt.Errorf("failed to auto-migrate webhook_logs table: %w", err)
	}
	return nil
}
