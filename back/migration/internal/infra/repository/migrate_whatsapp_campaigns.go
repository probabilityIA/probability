package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWhatsappCampaigns(ctx context.Context) error {
	conn := r.db.Conn(ctx)

	if err := conn.AutoMigrate(
		&models.WhatsappCampaign{},
		&models.WhatsappCampaignSend{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate whatsapp campaigns: %w", err)
	}

	statements := []string{
		`ALTER TABLE whatsapp_conversations
			ADD COLUMN IF NOT EXISTS campaign_id bigint`,
		`CREATE INDEX IF NOT EXISTS idx_whatsapp_campaign_id
			ON whatsapp_conversations (campaign_id)
			WHERE campaign_id IS NOT NULL`,
		`ALTER TABLE whatsapp_conversations
			DROP CONSTRAINT IF EXISTS chk_whatsapp_conversations_conversation_type`,
		`ALTER TABLE whatsapp_conversations
			ADD CONSTRAINT chk_whatsapp_conversations_conversation_type
			CHECK (conversation_type IN ('order','system_alert','inbound','campaign'))`,
		`CREATE INDEX IF NOT EXISTS idx_campaign_sends_pending
			ON whatsapp_campaign_sends (campaign_id, status)
			WHERE deleted_at IS NULL`,
		`CREATE INDEX IF NOT EXISTS idx_campaigns_due
			ON whatsapp_campaigns (scheduled_at)
			WHERE status IN ('scheduled','running') AND deleted_at IS NULL`,
	}

	for _, statement := range statements {
		if err := conn.Exec(statement).Error; err != nil {
			return fmt.Errorf("failed to apply campaign migration statement: %w", err)
		}
	}

	return nil
}
