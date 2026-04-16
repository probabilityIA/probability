package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
)

func (r *repository) Exists(ctx context.Context, integrationTypeID uint, originalStatus string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).Model(&models.OrderStatusMapping{}).
		Where("integration_type_id = ? AND original_status = ? AND is_active = ?", integrationTypeID, originalStatus, true).
		Count(&count).Error
	return count > 0, err
}
