package repository

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

func (r *repository) GetChannelStatusByID(ctx context.Context, id uint) (*entities.ChannelStatusInfo, error) {
	var model models.IntegrationChannelStatus

	err := r.db.Conn(ctx).Preload("IntegrationType").First(&model, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrChannelStatusNotFound
		}
		return nil, err
	}

	result := model.ToDomain()
	return &result, nil
}
