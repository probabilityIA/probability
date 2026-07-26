package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateNotificationConfigCODOnly(ctx context.Context) error {
	db := r.db.Conn(ctx)

	if err := db.Exec(`
ALTER TABLE business_notification_configs
DROP COLUMN IF EXISTS cod_only
`).Error; err != nil {
		return fmt.Errorf("drop cod_only from business_notification_configs: %w", err)
	}

	return nil
}
