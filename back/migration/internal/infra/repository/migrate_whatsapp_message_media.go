package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateWhatsappMessageMedia(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.WhatsAppMessageLog{}); err != nil {
		return fmt.Errorf("failed to auto-migrate whatsapp message media: %w", err)
	}
	return nil
}
