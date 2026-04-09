package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ListCustomerOrderItems(ctx context.Context, params dtos.ListCustomerOrderItemsParams) ([]entities.CustomerOrderItem, int64, error) {
	var modelsList []models.CustomerOrderItem
	var total int64

	query := r.db.Conn(ctx).Model(&models.CustomerOrderItem{}).
		Where("customer_id = ? AND business_id = ?", params.CustomerID, params.BusinessID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := paginationOffset(params.Page, params.PageSize)
	if err := query.Offset(offset).Limit(params.PageSize).Order("ordered_at DESC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entities.CustomerOrderItem, len(modelsList))
	for i, m := range modelsList {
		result[i] = mapCustomerOrderItemToEntity(&m)
	}
	return result, total, nil
}
