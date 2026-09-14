package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWhatsappConversationReads(ctx context.Context) error {
	conn := r.db.Conn(ctx)

	existed := conn.Migrator().HasTable(&models.WhatsAppConversationRead{})

	if err := conn.AutoMigrate(&models.WhatsAppConversationRead{}); err != nil {
		return fmt.Errorf("failed to auto-migrate whatsapp conversation reads: %w", err)
	}

	if existed {
		return nil
	}

	backfill := `
		INSERT INTO whatsapp_conversation_reads (business_id, phone_key, last_read_at, updated_at)
		SELECT DISTINCT business_id, regexp_replace(phone_number, '[^0-9]', '', 'g'), NOW(), NOW()
		FROM whatsapp_conversations
		WHERE regexp_replace(phone_number, '[^0-9]', '', 'g') <> ''
		ON CONFLICT (business_id, phone_key) DO NOTHING
	`
	if err := conn.Exec(backfill).Error; err != nil {
		return fmt.Errorf("failed to backfill whatsapp conversation reads: %w", err)
	}

	return nil
}
