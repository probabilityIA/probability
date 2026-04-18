package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateRackLevel(ctx context.Context, level *entities.WarehouseRackLevel) (*entities.WarehouseRackLevel, error) {
	model := &models.WarehouseRackLevel{
		RackID:     level.RackID,
		BusinessID: level.BusinessID,
		Code:       level.Code,
		Ordinal:    level.Ordinal,
		IsActive:   level.IsActive,
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return rackLevelModelToEntity(model), nil
}

func (r *Repository) GetRackLevelByID(ctx context.Context, businessID, levelID uint) (*entities.WarehouseRackLevel, error) {
	var model models.WarehouseRackLevel
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", levelID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrRackLevelNotFound
		}
		return nil, err
	}
	return rackLevelModelToEntity(&model), nil
}

func (r *Repository) ListRackLevels(ctx context.Context, params dtos.ListRackLevelsParams) ([]entities.WarehouseRackLevel, int64, error) {
	var modelsList []models.WarehouseRackLevel
	var total int64

	query := r.db.Conn(ctx).Model(&models.WarehouseRackLevel{}).
		Where("business_id = ? AND rack_id = ?", params.BusinessID, params.RackID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("ordinal ASC, code ASC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	levels := make([]entities.WarehouseRackLevel, len(modelsList))
	for i, m := range modelsList {
		levels[i] = *rackLevelModelToEntity(&m)
	}
	return levels, total, nil
}

func (r *Repository) UpdateRackLevel(ctx context.Context, level *entities.WarehouseRackLevel) (*entities.WarehouseRackLevel, error) {
	updates := map[string]any{
		"code":      level.Code,
		"ordinal":   level.Ordinal,
		"is_active": level.IsActive,
	}
	res := r.db.Conn(ctx).Model(&models.WarehouseRackLevel{}).
		Where("id = ? AND business_id = ?", level.ID, level.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrRackLevelNotFound
	}
	return r.GetRackLevelByID(ctx, level.BusinessID, level.ID)
}

func (r *Repository) DeleteRackLevel(ctx context.Context, businessID, levelID uint) error {
	res := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", levelID, businessID).
		Delete(&models.WarehouseRackLevel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrRackLevelNotFound
	}
	return nil
}

func (r *Repository) RackLevelExistsByCode(ctx context.Context, rackID uint, code string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.WarehouseRackLevel{}).
		Where("rack_id = ? AND code = ?", rackID, code)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
