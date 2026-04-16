package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
)

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&models.OrderStatusMapping{}, id).Error
}
