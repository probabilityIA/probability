package repository

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

func (r *repository) UpdateChannelStatus(ctx context.Context, id uint, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
	var model models.IntegrationChannelStatus

	err := r.db.Conn(ctx).First(&model, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrChannelStatusNotFound
		}
		return nil, err
	}

	model.Code = status.Code
	model.Name = status.Name
	model.Description = status.Description
	model.IsActive = status.IsActive
	model.DisplayOrder = status.DisplayOrder

	if err := r.db.Conn(ctx).Save(&model).Error; err != nil {
		return nil, err
	}

	// Preload IntegrationType for response
	r.db.Conn(ctx).Preload("IntegrationType").First(&model, model.ID)

	result := model.ToDomain()
	return &result, nil
}
