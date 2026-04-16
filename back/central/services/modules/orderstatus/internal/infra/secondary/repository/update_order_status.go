package repository

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

func (r *repository) UpdateOrderStatus(ctx context.Context, id uint, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
	var model models.OrderStatus

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrOrderStatusNotFound
		}
		return nil, err
	}

	model.Code = status.Code
	model.Name = status.Name
	model.Description = status.Description
	model.Category = status.Category
	model.Color = status.Color
	model.Priority = status.Priority
	model.IsActive = status.IsActive

	if err := r.db.Conn(ctx).Save(&model).Error; err != nil {
		return nil, err
	}

	result := model.ToDomain()
	return &result, nil
}
