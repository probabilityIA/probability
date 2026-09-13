package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWhatsappFlows(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(
		&models.WhatsappFlow{},
		&models.WhatsappTemplateFlow{},
		&models.WhatsappCampaign{},
		&models.ScheduledNotificationRule{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate whatsapp flows: %w", err)
	}

	const uniqueName = `
		CREATE UNIQUE INDEX IF NOT EXISTS idx_wa_flows_business_name_unique
		ON whatsapp_flows (business_id, lower(name))
		WHERE deleted_at IS NULL
	`

	if err := db.Exec(uniqueName).Error; err != nil {
		return fmt.Errorf("failed to create whatsapp flow unique name index: %w", err)
	}

	return nil
}
