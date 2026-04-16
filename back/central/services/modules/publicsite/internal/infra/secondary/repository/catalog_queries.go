package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) ListActiveProducts(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.PublicProduct, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Conn(ctx).Model(&models.Product{}).
		Where("business_id = ? AND is_active = true AND deleted_at IS NULL", businessID)

	if filters.Search != "" {
		like := "%" + filters.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Offset(filters.Offset()).
		Limit(filters.PageSize).
		Order("name ASC").
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return mappers.ProductsToEntities(products), total, nil
}

func (r *Repository) GetProductByID(ctx context.Context, businessID uint, productID string) (*entities.PublicProduct, error) {
	var product models.Product
	err := r.db.Conn(ctx).
		Where("id = ? AND business_id = ? AND is_active = true AND deleted_at IS NULL", productID, businessID).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrProductNotFound
		}
		return nil, err
	}
	return mappers.ProductToEntity(&product), nil
}

func (r *Repository) GetFeaturedProducts(ctx context.Context, businessID uint, limit int) ([]entities.PublicProduct, error) {
	var products []models.Product
	err := r.db.Conn(ctx).
		Where("business_id = ? AND is_active = true AND is_featured = true AND deleted_at IS NULL", businessID).
		Order("name ASC").
		Limit(limit).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return mappers.ProductsToEntities(products), nil
}
