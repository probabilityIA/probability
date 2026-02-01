package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// GetIntegrationCategoryByID obtiene una categoría de integración por su ID
func (r *Repository) GetIntegrationCategoryByID(ctx context.Context, id uint) (*domain.IntegrationCategory, error) {
	var categoryModel models.IntegrationCategory

	err := r.db.Conn(ctx).
		Where("id = ?", id).
		First(&categoryModel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Error(ctx).Uint("category_id", id).Msg("Integration category not found")
			return nil, fmt.Errorf("categoría de integración con ID %d no encontrada", id)
		}
		r.log.Error(ctx).Err(err).Uint("category_id", id).Msg("Error getting integration category by ID")
		return nil, err
	}

	// Mapear modelo GORM a entidad de dominio
	category := &domain.IntegrationCategory{
		ID:               categoryModel.ID,
		Code:             categoryModel.Code,
		Name:             categoryModel.Name,
		Description:      categoryModel.Description,
		Icon:             categoryModel.Icon,
		Color:            categoryModel.Color,
		DisplayOrder:     categoryModel.DisplayOrder,
		ParentCategoryID: categoryModel.ParentCategoryID,
		IsActive:         categoryModel.IsActive,
		IsVisible:        categoryModel.IsVisible,
		CreatedAt:        categoryModel.CreatedAt,
		UpdatedAt:        categoryModel.UpdatedAt,
	}

	return category, nil
}

// ListIntegrationCategories lista todas las categorías de integración activas y visibles
func (r *Repository) ListIntegrationCategories(ctx context.Context) ([]*domain.IntegrationCategory, error) {
	var categoriesModel []models.IntegrationCategory

	// Obtener categorías activas y visibles, ordenadas por display_order
	err := r.db.Conn(ctx).
		Where("is_active = ? AND is_visible = ?", true, true).
		Order("display_order ASC").
		Find(&categoriesModel).Error

	if err != nil {
		r.log.Error(ctx).Err(err).Msg("Error listing integration categories")
		return nil, err
	}

	// Mapear modelos GORM a entidades de dominio
	categories := make([]*domain.IntegrationCategory, 0, len(categoriesModel))
	for _, categoryModel := range categoriesModel {
		categories = append(categories, &domain.IntegrationCategory{
			ID:               categoryModel.ID,
			Code:             categoryModel.Code,
			Name:             categoryModel.Name,
			Description:      categoryModel.Description,
			Icon:             categoryModel.Icon,
			Color:            categoryModel.Color,
			DisplayOrder:     categoryModel.DisplayOrder,
			ParentCategoryID: categoryModel.ParentCategoryID,
			IsActive:         categoryModel.IsActive,
			IsVisible:        categoryModel.IsVisible,
			CreatedAt:        categoryModel.CreatedAt,
			UpdatedAt:        categoryModel.UpdatedAt,
		})
	}

	return categories, nil
}
