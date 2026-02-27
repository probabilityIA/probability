package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateLocation(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
	model := &models.WarehouseLocation{
		WarehouseID:   location.WarehouseID,
		Name:          location.Name,
		Code:          location.Code,
		Type:          location.Type,
		IsActive:      location.IsActive,
		IsFulfillment: location.IsFulfillment,
		Capacity:      location.Capacity,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	location.ID = model.ID
	location.CreatedAt = model.CreatedAt
	location.UpdatedAt = model.UpdatedAt
	return location, nil
}

func (r *Repository) GetLocationByID(ctx context.Context, warehouseID, locationID uint) (*entities.WarehouseLocation, error) {
	var model models.WarehouseLocation
	err := r.db.Conn(ctx).
		Where("id = ? AND warehouse_id = ?", locationID, warehouseID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrLocationNotFound
		}
		return nil, err
	}
	return locationModelToEntity(&model), nil
}

func (r *Repository) ListLocations(ctx context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
	var modelsList []models.WarehouseLocation

	err := r.db.Conn(ctx).
		Where("warehouse_id = ?", params.WarehouseID).
		Order("created_at ASC").
		Find(&modelsList).Error
	if err != nil {
		return nil, err
	}

	locations := make([]entities.WarehouseLocation, len(modelsList))
	for i, m := range modelsList {
		locations[i] = *locationModelToEntity(&m)
	}
	return locations, nil
}

func (r *Repository) UpdateLocation(ctx context.Context, location *entities.WarehouseLocation) (*entities.WarehouseLocation, error) {
	model := &models.WarehouseLocation{
		Model:         gorm.Model{ID: location.ID},
		WarehouseID:   location.WarehouseID,
		Name:          location.Name,
		Code:          location.Code,
		Type:          location.Type,
		IsActive:      location.IsActive,
		IsFulfillment: location.IsFulfillment,
		Capacity:      location.Capacity,
	}

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	location.UpdatedAt = model.UpdatedAt
	return location, nil
}

func (r *Repository) DeleteLocation(ctx context.Context, warehouseID, locationID uint) error {
	result := r.db.Conn(ctx).
		Where("id = ? AND warehouse_id = ?", locationID, warehouseID).
		Delete(&models.WarehouseLocation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrLocationNotFound
	}
	return nil
}

func (r *Repository) LocationExistsByCode(ctx context.Context, warehouseID uint, code string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.WarehouseLocation{}).
		Where("warehouse_id = ? AND code = ?", warehouseID, code)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
