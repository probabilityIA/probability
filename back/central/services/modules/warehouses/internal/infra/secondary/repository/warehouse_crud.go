package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) Create(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error) {
	model := &models.Warehouse{
		BusinessID:    warehouse.BusinessID,
		Name:          warehouse.Name,
		Code:          warehouse.Code,
		Address:       warehouse.Address,
		City:          warehouse.City,
		State:         warehouse.State,
		Country:       warehouse.Country,
		ZipCode:       warehouse.ZipCode,
		Phone:         warehouse.Phone,
		ContactName:   warehouse.ContactName,
		ContactEmail:  warehouse.ContactEmail,
		IsActive:      warehouse.IsActive,
		IsDefault:     warehouse.IsDefault,
		IsFulfillment: warehouse.IsFulfillment,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	warehouse.ID = model.ID
	warehouse.CreatedAt = model.CreatedAt
	warehouse.UpdatedAt = model.UpdatedAt
	return warehouse, nil
}

func (r *Repository) GetByID(ctx context.Context, businessID, warehouseID uint) (*entities.Warehouse, error) {
	var model models.Warehouse
	err := r.db.Conn(ctx).
		Preload("Locations", "deleted_at IS NULL").
		Where("id = ? AND business_id = ?", warehouseID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrWarehouseNotFound
		}
		return nil, err
	}
	return warehouseModelToEntity(&model), nil
}

func (r *Repository) List(ctx context.Context, params dtos.ListWarehousesParams) ([]entities.Warehouse, int64, error) {
	var modelsList []models.Warehouse
	var total int64

	query := r.db.Conn(ctx).Model(&models.Warehouse{}).
		Where("business_id = ?", params.BusinessID)

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}
	if params.IsFulfillment != nil {
		query = query.Where("is_fulfillment = ?", *params.IsFulfillment)
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR city ILIKE ?", like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("created_at DESC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	warehouses := make([]entities.Warehouse, len(modelsList))
	for i, m := range modelsList {
		warehouses[i] = *warehouseModelToEntity(&m)
	}
	return warehouses, total, nil
}

func (r *Repository) Update(ctx context.Context, warehouse *entities.Warehouse) (*entities.Warehouse, error) {
	model := &models.Warehouse{
		Model:         gorm.Model{ID: warehouse.ID},
		BusinessID:    warehouse.BusinessID,
		Name:          warehouse.Name,
		Code:          warehouse.Code,
		Address:       warehouse.Address,
		City:          warehouse.City,
		State:         warehouse.State,
		Country:       warehouse.Country,
		ZipCode:       warehouse.ZipCode,
		Phone:         warehouse.Phone,
		ContactName:   warehouse.ContactName,
		ContactEmail:  warehouse.ContactEmail,
		IsActive:      warehouse.IsActive,
		IsDefault:     warehouse.IsDefault,
		IsFulfillment: warehouse.IsFulfillment,
	}

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	warehouse.UpdatedAt = model.UpdatedAt
	return warehouse, nil
}

func (r *Repository) Delete(ctx context.Context, businessID, warehouseID uint) error {
	result := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", warehouseID, businessID).
		Delete(&models.Warehouse{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainerrors.ErrWarehouseNotFound
	}
	return nil
}

func (r *Repository) ExistsByCode(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.Warehouse{}).
		Where("business_id = ? AND code = ?", businessID, code)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *Repository) ClearDefault(ctx context.Context, businessID uint, excludeID uint) error {
	query := r.db.Conn(ctx).Model(&models.Warehouse{}).
		Where("business_id = ? AND is_default = ?", businessID, true)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	return query.Update("is_default", false).Error
}
