package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateLPN(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error) {
	m := mappers.LPNEntityToModel(lpn)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.LPNModelToEntity(m), nil
}

func (r *Repository) GetLPNByID(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
	var m models.LicensePlate
	err := r.db.Conn(ctx).Preload("Lines").Where("id = ? AND business_id = ?", id, businessID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrLPNNotFound
		}
		return nil, err
	}
	return mappers.LPNModelToEntity(&m), nil
}

func (r *Repository) GetLPNByCode(ctx context.Context, businessID uint, code string) (*entities.LicensePlate, error) {
	var m models.LicensePlate
	err := r.db.Conn(ctx).Preload("Lines").Where("business_id = ? AND code = ?", businessID, code).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrLPNNotFound
		}
		return nil, err
	}
	return mappers.LPNModelToEntity(&m), nil
}

func (r *Repository) ListLPNs(ctx context.Context, params dtos.ListLPNParams) ([]entities.LicensePlate, int64, error) {
	var ml []models.LicensePlate
	var total int64
	q := r.db.Conn(ctx).Model(&models.LicensePlate{}).Where("business_id = ?", params.BusinessID)
	if params.LpnType != "" {
		q = q.Where("lpn_type = ?", params.LpnType)
	}
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if params.LocationID != nil {
		q = q.Where("current_location_id = ?", *params.LocationID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(params.Offset()).Limit(params.PageSize).Order("id DESC").Find(&ml).Error; err != nil {
		return nil, 0, err
	}
	out := make([]entities.LicensePlate, len(ml))
	for i := range ml {
		out[i] = *mappers.LPNModelToEntity(&ml[i])
	}
	return out, total, nil
}

func (r *Repository) UpdateLPN(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error) {
	updates := map[string]any{
		"code":                lpn.Code,
		"lpn_type":            lpn.LpnType,
		"current_location_id": lpn.CurrentLocationID,
		"status":              lpn.Status,
	}
	res := r.db.Conn(ctx).Model(&models.LicensePlate{}).
		Where("id = ? AND business_id = ?", lpn.ID, lpn.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrLPNNotFound
	}
	return r.GetLPNByID(ctx, lpn.BusinessID, lpn.ID)
}

func (r *Repository) DeleteLPN(ctx context.Context, businessID, id uint) error {
	res := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).Delete(&models.LicensePlate{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrLPNNotFound
	}
	return nil
}

func (r *Repository) LPNExistsByCode(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
	var count int64
	q := r.db.Conn(ctx).Model(&models.LicensePlate{}).Where("business_id = ? AND code = ?", businessID, code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *Repository) AddLPNLine(ctx context.Context, line *entities.LicensePlateLine) (*entities.LicensePlateLine, error) {
	m := mappers.LPNLineEntityToModel(line)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.LPNLineModelToEntity(m), nil
}

func (r *Repository) ListLPNLines(ctx context.Context, lpnID uint) ([]entities.LicensePlateLine, error) {
	var ml []models.LicensePlateLine
	if err := r.db.Conn(ctx).Where("lpn_id = ?", lpnID).Find(&ml).Error; err != nil {
		return nil, err
	}
	out := make([]entities.LicensePlateLine, len(ml))
	for i := range ml {
		out[i] = *mappers.LPNLineModelToEntity(&ml[i])
	}
	return out, nil
}

func (r *Repository) DissolveLPN(ctx context.Context, businessID, id uint) error {
	res := r.db.Conn(ctx).Model(&models.LicensePlate{}).
		Where("id = ? AND business_id = ?", id, businessID).
		Updates(map[string]any{"status": "dissolved"})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrLPNNotFound
	}
	return nil
}
