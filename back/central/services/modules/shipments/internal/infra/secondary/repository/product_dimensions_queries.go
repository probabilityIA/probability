package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (r *Repository) GetProductDimensionsBySKUs(ctx context.Context, businessID uint, skus []string) (map[string]domain.ProductDimensions, error) {
	result := make(map[string]domain.ProductDimensions, len(skus))
	if len(skus) == 0 {
		return result, nil
	}

	var rows []struct {
		SKU    string
		Weight *float64
		Length *float64
		Width  *float64
		Height *float64
	}
	err := r.db.Conn(ctx).
		Table("products").
		Select("sku, weight, length, width, height").
		Where("business_id = ? AND sku IN ? AND deleted_at IS NULL", businessID, skus).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.SKU] = domain.ProductDimensions{
			Weight: row.Weight,
			Length: row.Length,
			Width:  row.Width,
			Height: row.Height,
		}
	}
	return result, nil
}
