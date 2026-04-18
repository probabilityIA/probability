package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateRack(ctx context.Context, rack *entities.WarehouseRack) (*entities.WarehouseRack, error) {
	model := &models.WarehouseRack{
		AisleID:     rack.AisleID,
		BusinessID:  rack.BusinessID,
		Code:        rack.Code,
		Name:        rack.Name,
		LevelsCount: rack.LevelsCount,
		IsActive:    rack.IsActive,
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return rackModelToEntity(model), nil
}

func (r *Repository) GetRackByID(ctx context.Context, businessID, rackID uint) (*entities.WarehouseRack, error) {
	var model models.WarehouseRack
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", rackID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrRackNotFound
		}
		return nil, err
	}
	return rackModelToEntity(&model), nil
}

func (r *Repository) ListRacks(ctx context.Context, params dtos.ListRacksParams) ([]entities.WarehouseRack, int64, error) {
	var modelsList []models.WarehouseRack
	var total int64

	query := r.db.Conn(ctx).Model(&models.WarehouseRack{}).
		Where("business_id = ? AND aisle_id = ?", params.BusinessID, params.AisleID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("code ASC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	racks := make([]entities.WarehouseRack, len(modelsList))
	for i, m := range modelsList {
		racks[i] = *rackModelToEntity(&m)
	}
	return racks, total, nil
}

func (r *Repository) UpdateRack(ctx context.Context, rack *entities.WarehouseRack) (*entities.WarehouseRack, error) {
	updates := map[string]any{
		"code":         rack.Code,
		"name":         rack.Name,
		"levels_count": rack.LevelsCount,
		"is_active":    rack.IsActive,
	}
	res := r.db.Conn(ctx).Model(&models.WarehouseRack{}).
		Where("id = ? AND business_id = ?", rack.ID, rack.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrRackNotFound
	}
	return r.GetRackByID(ctx, rack.BusinessID, rack.ID)
}

func (r *Repository) DeleteRack(ctx context.Context, businessID, rackID uint) error {
	res := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", rackID, businessID).
		Delete(&models.WarehouseRack{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrRackNotFound
	}
	return nil
}

func (r *Repository) RackExistsByCode(ctx context.Context, aisleID uint, code string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.WarehouseRack{}).
		Where("aisle_id = ? AND code = ?", aisleID, code)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
