package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migratePasswordResetTokens(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.PasswordResetToken{}); err != nil {
		return fmt.Errorf("failed to auto-migrate password_reset_tokens: %w", err)
	}
	return nil
}
