package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateBoldWebhookRawLog(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.BoldWebhookRawLog{}); err != nil {
		return fmt.Errorf("failed to auto-migrate bold_webhook_raw_logs table: %w", err)
	}
	return nil
}
