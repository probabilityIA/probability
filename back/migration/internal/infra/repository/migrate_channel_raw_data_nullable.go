package repository

import (
	"context"
	"fmt"
)

func (r *Repository) migrateChannelRawDataNullable(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(
		`ALTER TABLE order_channel_metadata ALTER COLUMN raw_data DROP NOT NULL`,
	).Error; err != nil {
		return fmt.Errorf("migrateChannelRawDataNullable: %w", err)
	}
	return nil
}
