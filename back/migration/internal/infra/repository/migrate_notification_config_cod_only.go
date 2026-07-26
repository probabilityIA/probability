package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateNotificationConfigCODOnly(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE business_notification_configs
ADD COLUMN IF NOT EXISTS cod_only BOOLEAN NOT NULL DEFAULT false
`).Error; err != nil {
		return fmt.Errorf("add cod_only to business_notification_configs: %w", err)
	}

	return nil
}
