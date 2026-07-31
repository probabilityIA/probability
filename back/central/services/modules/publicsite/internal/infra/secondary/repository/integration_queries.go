package repository

import (
	"context"

	"gorm.io/gorm"
)

func (r *Repository) IsIntegrationActive(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error) {
	var result struct {
		IsActive bool
	}

	err := r.db.Conn(ctx).
		Table("integrations").
		Select("is_active").
		Where("business_id = ? AND integration_type_id = ? AND deleted_at IS NULL", businessID, integrationTypeID).
		Limit(1).
		First(&result).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return result.IsActive, nil
}
