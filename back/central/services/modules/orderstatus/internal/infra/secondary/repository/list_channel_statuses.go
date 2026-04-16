package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
)

func (r *repository) ListChannelStatuses(ctx context.Context, integrationTypeID uint, isActive *bool) ([]entities.ChannelStatusInfo, error) {
	var modelsList []models.IntegrationChannelStatus

	query := r.db.Conn(ctx).Model(&models.IntegrationChannelStatus{}).
		Preload("IntegrationType").
		Where("integration_type_id = ?", integrationTypeID)

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Order("display_order ASC, code ASC").Find(&modelsList).Error; err != nil {
		return nil, err
	}

	result := make([]entities.ChannelStatusInfo, len(modelsList))
	for i, m := range modelsList {
		result[i] = m.ToDomain()
	}
	return result, nil
}
