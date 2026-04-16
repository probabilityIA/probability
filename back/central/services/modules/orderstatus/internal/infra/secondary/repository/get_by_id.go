package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

func (r *repository) GetByID(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
	var model models.OrderStatusMapping

	err := r.db.Conn(ctx).
		Preload("IntegrationType").
		Preload("OrderStatus").
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrMappingNotFound
		}
		return nil, err
	}

	domain := model.ToDomain()
	return &domain, nil
}
