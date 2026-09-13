package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWhatsappTemplateFlows(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.WhatsappTemplateFlow{}); err != nil {
		return fmt.Errorf("failed to auto-migrate whatsapp template flows: %w", err)
	}

	const uniqueIndex = `
		CREATE UNIQUE INDEX IF NOT EXISTS idx_wa_flow_source_button_unique
		ON whatsapp_template_flows (source_template_id, lower(button_text))
		WHERE deleted_at IS NULL
	`

	if err := db.Exec(uniqueIndex).Error; err != nil {
		return fmt.Errorf("failed to create whatsapp template flow unique index: %w", err)
	}

	return nil
}
