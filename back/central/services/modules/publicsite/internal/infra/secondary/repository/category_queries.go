package repository

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

func (r *Repository) GetHiddenCategories(ctx context.Context, businessID uint) ([]string, error) {
	var row struct {
		HiddenCategories []byte
	}
	err := r.db.Conn(ctx).
		Table("business_website_config").
		Select("hidden_categories").
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Limit(1).
		First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(row.HiddenCategories) == 0 {
		return nil, nil
	}
	var hidden []string
	if err := json.Unmarshal(row.HiddenCategories, &hidden); err != nil {
		return nil, nil
	}
	return hidden, nil
}

func (r *Repository) ListPublicCategories(ctx context.Context, businessID uint, hidden []string) ([]string, error) {
	query := r.db.Conn(ctx).
		Table("products").
		Select("category").
		Where("business_id = ? AND is_active = true AND deleted_at IS NULL AND category <> ''", businessID).
		Group("category").
		Order("category ASC")
	if len(hidden) > 0 {
		query = query.Where("category NOT IN (?)", hidden)
	}
	var categories []string
	if err := query.Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
