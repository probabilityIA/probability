package repository

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

func (r *repository) DeleteOrderStatus(ctx context.Context, id uint) error {
	// Verificar que no haya mapeos dependientes
	var count int64
	r.db.Conn(ctx).Model(&models.OrderStatusMapping{}).
		Where("order_status_id = ?", id).
		Count(&count)
	if count > 0 {
		return domainerrors.ErrOrderStatusHasMappings
	}

	var model models.OrderStatus
	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domainerrors.ErrOrderStatusNotFound
		}
		return err
	}

	return r.db.Conn(ctx).Delete(&model).Error
}
