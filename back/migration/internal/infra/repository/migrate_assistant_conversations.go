package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateAssistantConversations(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.AssistantConversation{}, &models.AssistantMessage{}); err != nil {
		return fmt.Errorf("failed to auto-migrate assistant conversations: %w", err)
	}
	return nil
}
