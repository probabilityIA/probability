package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateScheduledNotifications(ctx context.Context) error {
	conn := r.db.Conn(ctx)

	if err := conn.AutoMigrate(&models.WhatsappTemplate{}); err != nil {
		return fmt.Errorf("failed to auto-migrate whatsapp_templates: %w", err)
	}

	if err := conn.AutoMigrate(&models.Client{}); err != nil {
		return fmt.Errorf("failed to add marketing opt-in columns to client: %w", err)
	}

	if err := conn.AutoMigrate(
		&models.ScheduledNotificationRule{},
		&models.ScheduledNotificationRun{},
		&models.ScheduledNotificationSend{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate scheduled notifications: %w", err)
	}

	statements := []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_whatsapp_templates_meta_id
			ON whatsapp_templates (meta_template_id)
			WHERE meta_template_id <> '' AND deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_scheduled_rules_due
			ON scheduled_notification_rules (next_run_at)
			WHERE enabled = true AND deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_sched_send_rule_client_recent
			ON scheduled_notification_sends (rule_id, client_id, queued_at DESC)`,
	}

	for _, statement := range statements {
		if err := conn.Exec(statement).Error; err != nil {
			return fmt.Errorf("failed to create scheduled notification index: %w", err)
		}
	}

	return nil
}
