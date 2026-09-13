package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWhatsappTemplateHeaderMedia(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.AutoMigrate(&models.WhatsappTemplate{}); err != nil {
		return fmt.Errorf("failed to auto-migrate whatsapp template header media: %w", err)
	}

	if err := db.Exec(`UPDATE whatsapp_templates SET header_type = 'TEXT' WHERE header_type IS NULL OR header_type = ''`).Error; err != nil {
		return fmt.Errorf("failed to backfill whatsapp template header type: %w", err)
	}

	return nil
}
