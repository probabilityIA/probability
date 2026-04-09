package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ListCustomerProducts(ctx context.Context, params dtos.ListCustomerProductsParams) ([]entities.CustomerProductHistory, int64, error) {
	var modelsList []models.CustomerProductHistory
	var total int64

	query := r.db.Conn(ctx).Model(&models.CustomerProductHistory{}).
		Where("customer_id = ? AND business_id = ?", params.CustomerID, params.BusinessID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := resolveProductSortOrder(params.SortBy)
	offset := paginationOffset(params.Page, params.PageSize)
	if err := query.Offset(offset).Limit(params.PageSize).Order(orderClause).Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entities.CustomerProductHistory, len(modelsList))
	for i, m := range modelsList {
		result[i] = mapCustomerProductToEntity(&m)
	}
	return result, total, nil
}

func resolveProductSortOrder(sortBy string) string {
	switch sortBy {
	case "total_spent":
		return "total_spent DESC"
	case "last_ordered_at":
		return "last_ordered_at DESC"
	default:
		return "times_ordered DESC"
	}
}
