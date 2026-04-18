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

func (r *Repository) ListUoMs(ctx context.Context) ([]entities.UnitOfMeasure, error) {
	var modelsList []models.UnitOfMeasure
	if err := r.db.Conn(ctx).Where("is_active = ?", true).Order("code ASC").Find(&modelsList).Error; err != nil {
		return nil, err
	}
	uoms := make([]entities.UnitOfMeasure, len(modelsList))
	for i := range modelsList {
		uoms[i] = *mappers.UoMModelToEntity(&modelsList[i])
	}
	return uoms, nil
}

func (r *Repository) GetUoMByCode(ctx context.Context, code string) (*entities.UnitOfMeasure, error) {
	var m models.UnitOfMeasure
	err := r.db.Conn(ctx).Where("code = ?", code).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrUomNotFound
		}
		return nil, err
	}
	return mappers.UoMModelToEntity(&m), nil
}

func (r *Repository) GetUoMByID(ctx context.Context, uomID uint) (*entities.UnitOfMeasure, error) {
	var m models.UnitOfMeasure
	err := r.db.Conn(ctx).First(&m, uomID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrUomNotFound
		}
		return nil, err
	}
	return mappers.UoMModelToEntity(&m), nil
}

func (r *Repository) CreateProductUoM(ctx context.Context, pu *entities.ProductUoM) (*entities.ProductUoM, error) {
	model := mappers.ProductUoMEntityToModel(pu)
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	var saved models.ProductUoM
	r.db.Conn(ctx).Preload("Uom").First(&saved, model.ID)
	return mappers.ProductUoMModelToEntity(&saved), nil
}

func (r *Repository) ListProductUoMs(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
	var modelsList []models.ProductUoM
	if err := r.db.Conn(ctx).Preload("Uom").
		Where("business_id = ? AND product_id = ?", params.BusinessID, params.ProductID).
		Order("is_base DESC, id ASC").
		Find(&modelsList).Error; err != nil {
		return nil, err
	}
	result := make([]entities.ProductUoM, len(modelsList))
	for i := range modelsList {
		result[i] = *mappers.ProductUoMModelToEntity(&modelsList[i])
	}
	return result, nil
}

func (r *Repository) DeleteProductUoM(ctx context.Context, businessID, id uint) error {
	res := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).Delete(&models.ProductUoM{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrProductUoMNotFound
	}
	return nil
}

func (r *Repository) GetBaseProductUoM(ctx context.Context, businessID uint, productID string) (*entities.ProductUoM, error) {
	var m models.ProductUoM
	err := r.db.Conn(ctx).Preload("Uom").
		Where("business_id = ? AND product_id = ? AND is_base = true", businessID, productID).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrProductUoMNotFound
		}
		return nil, err
	}
	return mappers.ProductUoMModelToEntity(&m), nil
}
