package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ListChannelPaymentMethods(ctx context.Context, integrationType *string, isActive *bool) ([]entities.ChannelPaymentMethod, error) {
	var rows []models.ChannelPaymentMethod

	query := r.db.Conn(ctx).Model(&models.ChannelPaymentMethod{})
	if integrationType != nil {
		query = query.Where("integration_type = ?", *integrationType)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Order("integration_type, display_order").Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]entities.ChannelPaymentMethod, len(rows))
	for i, row := range rows {
		result[i] = entities.ChannelPaymentMethod{
			ID:              row.ID,
			IntegrationType: row.IntegrationType,
			Code:            row.Code,
			Name:            row.Name,
			Description:     row.Description,
			IsActive:        row.IsActive,
			DisplayOrder:    row.DisplayOrder,
		}
	}
	return result, nil
}
