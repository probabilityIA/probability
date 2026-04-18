package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"gorm.io/gorm"
)

func (r *Repository) GetLocationCapacity(ctx context.Context, locationID uint) (*ports.LocationCapacityInfo, error) {
	var result struct {
		ID           uint
		MaxWeightKg  *float64
		MaxVolumeCm3 *float64
	}
	err := r.db.Conn(ctx).
		Table("warehouse_locations").
		Select("id, max_weight_kg, max_volume_cm3").
		Where("id = ? AND deleted_at IS NULL", locationID).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &ports.LocationCapacityInfo{
		ID:           result.ID,
		MaxWeightKg:  result.MaxWeightKg,
		MaxVolumeCm3: result.MaxVolumeCm3,
	}, nil
}

func (r *Repository) GetProductDimensions(ctx context.Context, productID string, businessID uint) (*ports.ProductDimensions, error) {
	var result struct {
		ID            string
		Weight        float64
		WeightUnit    string
		Length        float64
		Width         float64
		Height        float64
		DimensionUnit string
	}
	err := r.db.Conn(ctx).
		Table("products").
		Select("id, weight, weight_unit, length, width, height, dimension_unit").
		Where("id = ? AND business_id = ? AND deleted_at IS NULL", productID, businessID).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &ports.ProductDimensions{
		ID:      result.ID,
		Weight:  result.Weight,
		WeightU: result.WeightUnit,
		Length:  result.Length,
		Width:   result.Width,
		Height:  result.Height,
		DimU:    result.DimensionUnit,
	}, nil
}

func (r *Repository) GetLocationOccupiedQty(ctx context.Context, locationID uint) (int, error) {
	var total struct {
		Sum int
	}
	err := r.db.Conn(ctx).
		Table("inventory_levels").
		Select("COALESCE(SUM(quantity), 0) as sum").
		Where("location_id = ? AND deleted_at IS NULL", locationID).
		Scan(&total).Error
	return total.Sum, err
}
