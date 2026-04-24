package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateZone(ctx context.Context, zone *entities.WarehouseZone) (*entities.WarehouseZone, error) {
	model := &models.WarehouseZone{
		WarehouseID: zone.WarehouseID,
		BusinessID:  zone.BusinessID,
		Code:        zone.Code,
		Name:        zone.Name,
		Purpose:     zone.Purpose,
		IsActive:    zone.IsActive,
		ColorHex:    zone.ColorHex,
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return zoneModelToEntity(model), nil
}

func (r *Repository) GetZoneByID(ctx context.Context, businessID, zoneID uint) (*entities.WarehouseZone, error) {
	var model models.WarehouseZone
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", zoneID, businessID).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrZoneNotFound
		}
		return nil, err
	}
	return zoneModelToEntity(&model), nil
}

func (r *Repository) ListZones(ctx context.Context, params dtos.ListZonesParams) ([]entities.WarehouseZone, int64, error) {
	var modelsList []models.WarehouseZone
	var total int64

	query := r.db.Conn(ctx).Model(&models.WarehouseZone{}).
		Where("business_id = ? AND warehouse_id = ?", params.BusinessID, params.WarehouseID)

	if params.ActiveOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("code ASC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	zones := make([]entities.WarehouseZone, len(modelsList))
	for i, m := range modelsList {
		zones[i] = *zoneModelToEntity(&m)
	}
	return zones, total, nil
}

func (r *Repository) UpdateZone(ctx context.Context, zone *entities.WarehouseZone) (*entities.WarehouseZone, error) {
	updates := map[string]any{
		"code":      zone.Code,
		"name":      zone.Name,
		"purpose":   zone.Purpose,
		"is_active": zone.IsActive,
		"color_hex": zone.ColorHex,
	}
	res := r.db.Conn(ctx).Model(&models.WarehouseZone{}).
		Where("id = ? AND business_id = ?", zone.ID, zone.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrZoneNotFound
	}
	return r.GetZoneByID(ctx, zone.BusinessID, zone.ID)
}

func (r *Repository) DeleteZone(ctx context.Context, businessID, zoneID uint) error {
	res := r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", zoneID, businessID).
		Delete(&models.WarehouseZone{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrZoneNotFound
	}
	return nil
}

func (r *Repository) ZoneExistsByCode(ctx context.Context, warehouseID uint, code string, excludeID *uint) (bool, error) {
	var count int64
	query := r.db.Conn(ctx).Model(&models.WarehouseZone{}).
		Where("warehouse_id = ? AND code = ?", warehouseID, code)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
