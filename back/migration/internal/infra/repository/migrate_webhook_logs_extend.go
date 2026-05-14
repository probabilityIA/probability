package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWebhookLogsExtend(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.WebhookLog{}); err != nil {
		return fmt.Errorf("failed to extend webhook_logs table: %w", err)
	}
	if err := r.db.Conn(ctx).Exec("DROP TABLE IF EXISTS bold_webhook_raw_logs").Error; err != nil {
		return fmt.Errorf("failed to drop bold_webhook_raw_logs: %w", err)
	}
	return nil
}
