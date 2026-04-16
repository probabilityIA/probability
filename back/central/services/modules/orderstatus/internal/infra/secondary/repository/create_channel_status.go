package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
)

func (r *repository) CreateChannelStatus(ctx context.Context, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
	model := models.IntegrationChannelStatus{
		IntegrationTypeID: status.IntegrationTypeID,
		Code:              status.Code,
		Name:              status.Name,
		Description:       status.Description,
		IsActive:          status.IsActive,
		DisplayOrder:      status.DisplayOrder,
	}

	if err := r.db.Conn(ctx).Create(&model).Error; err != nil {
		return nil, err
	}

	// Preload IntegrationType for response
	r.db.Conn(ctx).Preload("IntegrationType").First(&model, model.ID)

	result := model.ToDomain()
	return &result, nil
}
