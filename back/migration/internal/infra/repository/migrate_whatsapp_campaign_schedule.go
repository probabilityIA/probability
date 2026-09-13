package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWhatsappCampaignSchedule(ctx context.Context) error {
	conn := r.db.Conn(ctx)

	if err := conn.Exec(`ALTER TABLE whatsapp_campaign_sends ADD COLUMN IF NOT EXISTS round bigint NOT NULL DEFAULT 1`).Error; err != nil {
		return fmt.Errorf("failed to add campaign send round: %w", err)
	}

	if err := conn.Exec(`DROP INDEX IF EXISTS idx_campaign_send_campaign_client`).Error; err != nil {
		return fmt.Errorf("failed to drop old campaign send unique index: %w", err)
	}

	if err := conn.AutoMigrate(
		&models.WhatsappCampaign{},
		&models.WhatsappCampaignSend{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate whatsapp campaign schedule: %w", err)
	}

	return nil
}
