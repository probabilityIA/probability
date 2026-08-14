package repository

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"gorm.io/datatypes"
)

func (r *Repository) GetWarehouseShippingConfig(ctx context.Context, warehouseID uint) (*domain.ShippingPackageConfig, error) {
	var row struct {
		Metadata datatypes.JSON
	}
	err := r.db.Conn(ctx).
		Table("warehouses").
		Select("metadata").
		Where("id = ?", warehouseID).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if len(row.Metadata) == 0 {
		return nil, nil
	}

	var raw struct {
		ShippingPackageStrategy string `json:"shipping_package_strategy"`
		StandardBoxes           []struct {
			Name     string   `json:"name"`
			Weight   *float64 `json:"weight"`
			Length   *float64 `json:"length"`
			Width    *float64 `json:"width"`
			Height   *float64 `json:"height"`
			MaxItems int      `json:"max_items"`
		} `json:"standard_boxes"`
	}
	if err := json.Unmarshal(row.Metadata, &raw); err != nil {
		return nil, nil
	}
	if raw.ShippingPackageStrategy == "" {
		return nil, nil
	}

	boxes := make([]domain.StandardBox, 0, len(raw.StandardBoxes))
	for _, b := range raw.StandardBoxes {
		boxes = append(boxes, domain.StandardBox{
			Name:     b.Name,
			Weight:   b.Weight,
			Length:   b.Length,
			Width:    b.Width,
			Height:   b.Height,
			MaxItems: b.MaxItems,
		})
	}

	return &domain.ShippingPackageConfig{
		Strategy:      raw.ShippingPackageStrategy,
		StandardBoxes: boxes,
	}, nil
}

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
