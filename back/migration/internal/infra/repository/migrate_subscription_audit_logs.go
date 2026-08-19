package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateSubscriptionAuditLogs(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.SubscriptionAuditLog{}); err != nil {
		return fmt.Errorf("failed to auto-migrate subscription_audit_logs: %w", err)
	}

	return nil
}
