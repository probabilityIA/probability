package repository

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

func (r *repository) GetOrderStatusByID(ctx context.Context, id uint) (*entities.OrderStatusInfo, error) {
	var model models.OrderStatus

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrOrderStatusNotFound
		}
		return nil, err
	}

	result := model.ToDomain()
	return &result, nil
}
