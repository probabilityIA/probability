package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// UpdateLastSync actualiza el timestamp de última sincronización de una integración
func (r *Repository) UpdateLastSync(ctx context.Context, id uint, lastSync time.Time) error {
	ctx = log.WithFunctionCtx(ctx, "UpdateLastSync")

	result := r.db.Conn(ctx).Model(&domain.Integration{}).
		Where("id = ?", id).
		Update("last_sync", lastSync)

	if result.Error != nil {
		r.log.Error(ctx).Err(result.Error).Uint("id", id).Msg("Error updating last_sync")
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Warn(ctx).Uint("id", id).Msg("No integration found to update last_sync")
	}

	r.log.Info(ctx).Uint("id", id).Time("last_sync", lastSync).Msg("Last sync updated successfully")
	return nil
}
