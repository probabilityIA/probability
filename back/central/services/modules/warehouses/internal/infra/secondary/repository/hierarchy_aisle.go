package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateAisle(ctx context.Context, aisle *entities.WarehouseAisle) (*entities.WarehouseAisle, error) {
	model := &models.WarehouseAisle{
		ZoneID:     aisle.ZoneID,
		BusinessID: aisle.BusinessID,
		Code:       aisle.Code,
		Name:       aisle.Name,
		IsActive:   aisle.IsActive,
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return aisleModelToEntity(model), nil
}

func (r *Repository) GetAisleByID(ctx context.Context, businessID, aisleID uint) (*entities.WarehouseAisle, error) {
	var model models.WarehouseAisle
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", aisleID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrAisleNotFound
		}
		return nil, err
	}
	return aisleModelToEntity(&model), nil
}

func (r *Repository) ListAisles(ctx context.Context, params dtos.ListAislesParams) ([]entities.WarehouseAisle, int64, error) {
	var modelsList []models.WarehouseAisle
	var total int64

	query := r.db.Conn(ctx).Model(&models.WarehouseAisle{}).
		Where("business_id = ? AND zone_id = ?", params.BusinessID, params.ZoneID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("code ASC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	aisles := make([]entities.WarehouseAisle, len(modelsList))
	for i, m := range modelsList {
		aisles[i] = *aisleModelToEntity(&m)
	}
	return aisles, total, nil
}

func (r *Repository) UpdateAisle(ctx context.Context, aisle *entities.WarehouseAisle) (*entities.WarehouseAisle, error) {
	updates := map[string]any{
		"code":      aisle.Code,
		"name":      aisle.Name,
		"is_active": aisle.IsActive,
	}
	res := r.db.Conn(ctx).Model(&models.WarehouseAisle{}).
		Where("id = ? AND business_id = ?", aisle.ID, aisle.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrAisleNotFound
	}
	return r.GetAisleByID(ctx, aisle.BusinessID, aisle.ID)
}

func (r *Repository) DeleteAisle(ctx context.Context, businessID, aisleID uint) error {
	res := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", aisleID, businessID).
		Delete(&models.WarehouseAisle{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrAisleNotFound
	}
	return nil
}

func (r *Repository) AisleExistsByCode(ctx context.Context, zoneID uint, code string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.WarehouseAisle{}).
		Where("zone_id = ? AND code = ?", zoneID, code)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
