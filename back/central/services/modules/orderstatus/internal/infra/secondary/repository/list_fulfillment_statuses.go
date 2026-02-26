package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *repository) ListFulfillmentStatuses(ctx context.Context, isActive *bool) ([]entities.FulfillmentStatusInfo, error) {
	var statuses []models.FulfillmentStatus

	query := r.db.Conn(ctx).Model(&models.FulfillmentStatus{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("code ASC").Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.FulfillmentStatusInfo, len(statuses))
	for i, s := range statuses {
		result[i] = entities.FulfillmentStatusInfo{
			ID:          s.ID,
			Code:        s.Code,
			Name:        s.Name,
			Description: s.Description,
			Category:    s.Category,
			Color:       s.Color,
		}
	}

	return result, nil
}
