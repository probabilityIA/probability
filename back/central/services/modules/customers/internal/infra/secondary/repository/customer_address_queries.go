package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ListCustomerAddresses(ctx context.Context, params dtos.ListCustomerAddressesParams) ([]entities.CustomerAddress, int64, error) {
	var modelsList []models.CustomerAddress
	var total int64

	query := r.db.Conn(ctx).Model(&models.CustomerAddress{}).
		Where("customer_id = ? AND business_id = ?", params.CustomerID, params.BusinessID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := paginationOffset(params.Page, params.PageSize)
	if err := query.Offset(offset).Limit(params.PageSize).Order("times_used DESC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entities.CustomerAddress, len(modelsList))
	for i, m := range modelsList {
		result[i] = mapCustomerAddressToEntity(&m)
	}
	return result, total, nil
}
