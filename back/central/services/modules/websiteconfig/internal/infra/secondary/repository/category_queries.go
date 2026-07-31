package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
)

func (r *Repository) ListProductCategories(ctx context.Context, businessID uint) ([]entities.ProductCategory, error) {
	var rows []struct {
		Category string
		Total    int64
	}
	err := r.db.Conn(ctx).
		Table("products").
		Select("category, COUNT(*) AS total").
		Where("business_id = ? AND is_active = true AND deleted_at IS NULL AND category <> ''", businessID).
		Group("category").
		Order("category ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	categories := make([]entities.ProductCategory, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, entities.ProductCategory{Name: row.Category, ProductCount: row.Total})
	}
	return categories, nil
}
