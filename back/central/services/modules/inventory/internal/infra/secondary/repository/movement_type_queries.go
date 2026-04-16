package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func movementTypeModelToEntity(m *models.StockMovementType) *entities.StockMovementType {
	return &entities.StockMovementType{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		IsActive:    m.IsActive,
		Direction:   m.Direction,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *Repository) ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error) {
	var modelsList []models.StockMovementType
	var total int64

	query := r.db.Conn(ctx).Model(&models.StockMovementType{})

	if params.ActiveOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(params.Offset()).Limit(params.PageSize).Order("id ASC").Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entities.StockMovementType, len(modelsList))
	for i, m := range modelsList {
		result[i] = *movementTypeModelToEntity(&m)
	}
	return result, total, nil
}

func (r *Repository) GetMovementTypeByID(ctx context.Context, id uint) (*entities.StockMovementType, error) {
	var model models.StockMovementType
	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrMovementTypeNotFound
		}
		return nil, err
	}
	return movementTypeModelToEntity(&model), nil
}

func (r *Repository) GetMovementTypeIDByCode(ctx context.Context, code string) (uint, error) {
	var result struct {
		ID uint
	}
	err := r.db.Conn(ctx).
		Model(&models.StockMovementType{}).
		Select("id").
		Where("code = ? AND is_active = ?", code, true).
		First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, domainerrors.ErrMovementTypeNotFound
		}
		return 0, err
	}
	return result.ID, nil
}

func (r *Repository) CreateMovementType(ctx context.Context, movType *entities.StockMovementType) (*entities.StockMovementType, error) {
	// Verificar que el code no exista
	var count int64
	r.db.Conn(ctx).Model(&models.StockMovementType{}).Where("code = ?", movType.Code).Count(&count)
	if count > 0 {
		return nil, domainerrors.ErrMovementTypeCodeExists
	}

	model := &models.StockMovementType{
		Code:        movType.Code,
		Name:        movType.Name,
		Description: movType.Description,
		IsActive:    movType.IsActive,
		Direction:   movType.Direction,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return movementTypeModelToEntity(model), nil
}

func (r *Repository) UpdateMovementType(ctx context.Context, movType *entities.StockMovementType) error {
	return r.db.Conn(ctx).Model(&models.StockMovementType{}).
		Where("id = ?", movType.ID).
		Updates(map[string]interface{}{
			"name":        movType.Name,
			"description": movType.Description,
			"is_active":   movType.IsActive,
			"direction":   movType.Direction,
		}).Error
}

func (r *Repository) DeleteMovementType(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&models.StockMovementType{}, id).Error
}
