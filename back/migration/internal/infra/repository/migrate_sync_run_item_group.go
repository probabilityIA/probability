package repository

import "context"

func (r *Repository) migrateSyncRunItemGroup(ctx context.Context) error {
	if err := r.db.Conn(ctx).Exec(
		`ALTER TABLE integration_sync_run_items ADD COLUMN IF NOT EXISTS group_ref varchar(64)`,
	).Error; err != nil {
		return err
	}
	if err := r.db.Conn(ctx).Exec(
		`ALTER TABLE integration_sync_run_items ADD COLUMN IF NOT EXISTS group_label varchar(300)`,
	).Error; err != nil {
		return err
	}
	return r.db.Conn(ctx).Exec(
		`CREATE INDEX IF NOT EXISTS idx_integration_sync_run_items_group_ref ON integration_sync_run_items(group_ref)`,
	).Error
}
