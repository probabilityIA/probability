package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateWhatsappInboundConversationType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE whatsapp_conversations DROP CONSTRAINT IF EXISTS chk_whatsapp_conversations_conversation_type
`).Error; err != nil {
		return fmt.Errorf("drop existing conversation_type check: %w", err)
	}

	if err := db.Exec(`
ALTER TABLE whatsapp_conversations ADD CONSTRAINT chk_whatsapp_conversations_conversation_type
    CHECK (conversation_type IN ('order', 'system_alert', 'inbound'))
`).Error; err != nil {
		return fmt.Errorf("add conversation_type check with inbound: %w", err)
	}

	return nil
}
