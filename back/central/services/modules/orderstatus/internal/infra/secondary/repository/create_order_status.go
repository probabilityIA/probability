package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
)

func (r *repository) CreateOrderStatus(ctx context.Context, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
	model := models.OrderStatus{
		Code:        status.Code,
		Name:        status.Name,
		Description: status.Description,
		Category:    status.Category,
		Color:       status.Color,
		Priority:    status.Priority,
		IsActive:    status.IsActive,
	}

	if err := r.db.Conn(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	result := model.ToDomain()
	return &result, nil
}
