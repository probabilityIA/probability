package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
)

func (r *Repository) CheckStockForOrder(ctx context.Context, businessID uint, warehouseID *uint, items []dtos.StockCheckItem) ([]dtos.StockCheckResult, error) {
	if len(items) == 0 {
		return nil, nil
	}

	results := make([]dtos.StockCheckResult, 0, len(items))

	for _, it := range items {
		if it.ProductID == "" || it.Quantity <= 0 {
			continue
		}

		query := r.db.Conn(ctx).
			Table("inventory_levels").
			Select("COALESCE(SUM(available_qty), 0) AS available").
			Where("product_id = ?", it.ProductID).
			Where("business_id = ?", businessID).
			Where("deleted_at IS NULL")

		if warehouseID != nil && *warehouseID > 0 {
			query = query.Where("warehouse_id = ?", *warehouseID)
		}

		var row struct{ Available int64 }
		if err := query.Scan(&row).Error; err != nil {
			return nil, err
		}

		results = append(results, dtos.StockCheckResult{
			ProductID:  it.ProductID,
			ProductSKU: it.ProductSKU,
			Required:   it.Quantity,
			Available:  int(row.Available),
			Sufficient: int(row.Available) >= it.Quantity,
		})
	}

	return results, nil
}
