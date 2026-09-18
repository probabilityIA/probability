package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type warehouseOriginRow struct {
	ID        uint
	Name      string
	Address   string
	Street    string
	City      string
	Latitude  *float64
	Longitude *float64
}

func (row warehouseOriginRow) toEntity() *entities.OriginWarehouse {
	if row.ID == 0 {
		return nil
	}
	address := row.Address
	if address == "" {
		address = row.Street
	}
	return &entities.OriginWarehouse{
		ID:      row.ID,
		Name:    row.Name,
		Address: address,
		City:    row.City,
		Lat:     row.Latitude,
		Lng:     row.Longitude,
	}
}

func (r *Repository) GetWarehouseOrigin(ctx context.Context, businessID, warehouseID uint) (*entities.OriginWarehouse, error) {
	var row warehouseOriginRow
	err := r.db.Conn(ctx).
		Model(&models.Warehouse{}).
		Select("id, name, address, street, city, latitude, longitude").
		Where("id = ? AND business_id = ?", warehouseID, businessID).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return row.toEntity(), nil
}

func (r *Repository) GetDefaultWarehouseOrigin(ctx context.Context, businessID uint) (*entities.OriginWarehouse, error) {
	var row warehouseOriginRow
	err := r.db.Conn(ctx).
		Model(&models.Warehouse{}).
		Select("id, name, address, street, city, latitude, longitude").
		Where("business_id = ? AND is_active = true", businessID).
		Order("is_default DESC").
		Order("id ASC").
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return row.toEntity(), nil
}

func (r *Repository) GetOrdersWarehouseIDs(ctx context.Context, businessID uint, orderIDs []string) ([]uint, error) {
	if len(orderIDs) == 0 {
		return nil, nil
	}
	var ids []uint
	err := r.db.Conn(ctx).
		Model(&models.Order{}).
		Distinct("warehouse_id").
		Where("business_id = ? AND id IN ? AND warehouse_id IS NOT NULL", businessID, orderIDs).
		Pluck("warehouse_id", &ids).Error
	return ids, err
}
