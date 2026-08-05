package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateWhatsappConversationType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`ALTER TABLE whatsapp_conversations ADD COLUMN IF NOT EXISTS conversation_type varchar(20) NOT NULL DEFAULT 'order'`).Error; err != nil {
		return fmt.Errorf("add conversation_type to whatsapp_conversations: %w", err)
	}

	if err := db.Exec(`ALTER TABLE whatsapp_conversations ALTER COLUMN order_number DROP NOT NULL`).Error; err != nil {
		return fmt.Errorf("drop not null on order_number: %w", err)
	}

	if err := db.Exec(`
ALTER TABLE whatsapp_conversations DROP CONSTRAINT IF EXISTS chk_whatsapp_conversations_conversation_type
`).Error; err != nil {
		return fmt.Errorf("drop existing conversation_type check: %w", err)
	}

	if err := db.Exec(`
ALTER TABLE whatsapp_conversations ADD CONSTRAINT chk_whatsapp_conversations_conversation_type
    CHECK (conversation_type IN ('order', 'system_alert'))
`).Error; err != nil {
		return fmt.Errorf("add conversation_type check: %w", err)
	}

	return nil
}
