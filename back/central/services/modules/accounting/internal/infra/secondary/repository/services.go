package repository

import (
	"context"
	stderrors "errors"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func serviceToEntity(m *models.AccountingService) *entities.Service {
	return &entities.Service{
		ID:         m.ID,
		Code:       m.Code,
		Name:       m.Name,
		UnspscCode: m.UnspscCode,
		UnitPrice:  m.UnitPrice,
		IsActive:   m.IsActive,
	}
}

func (r *Repository) ListServices(ctx context.Context) ([]entities.Service, error) {
	var serviceModels []models.AccountingService
	if err := r.db.Conn(ctx).Order("code ASC").Find(&serviceModels).Error; err != nil {
		return nil, err
	}
	out := make([]entities.Service, len(serviceModels))
	for i := range serviceModels {
		out[i] = *serviceToEntity(&serviceModels[i])
	}
	return out, nil
}

func (r *Repository) GetServiceByID(ctx context.Context, id uint) (*entities.Service, error) {
	var model models.AccountingService
	err := r.db.Conn(ctx).First(&model, id).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrServiceNotFound
		}
		return nil, err
	}
	return serviceToEntity(&model), nil
}

func (r *Repository) ServiceCodeExists(ctx context.Context, code string, excludeID *uint) (bool, error) {
	var count int64
	q := r.db.Conn(ctx).Model(&models.AccountingService{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateService(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error) {
	model := &models.AccountingService{
		Code:       dto.Code,
		Name:       dto.Name,
		UnspscCode: dto.UnspscCode,
		UnitPrice:  dto.UnitPrice,
		IsActive:   dto.IsActive,
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return serviceToEntity(model), nil
}

func (r *Repository) UpdateService(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error) {
	updates := map[string]interface{}{
		"code":        dto.Code,
		"name":        dto.Name,
		"unspsc_code": dto.UnspscCode,
		"unit_price":  dto.UnitPrice,
		"is_active":   dto.IsActive,
	}
	res := r.db.Conn(ctx).Model(&models.AccountingService{}).Where("id = ?", dto.ID).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrServiceNotFound
	}
	return r.GetServiceByID(ctx, dto.ID)
}

func (r *Repository) DeleteService(ctx context.Context, id uint) error {
	res := r.db.Conn(ctx).Delete(&models.AccountingService{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrServiceNotFound
	}
	return nil
}
