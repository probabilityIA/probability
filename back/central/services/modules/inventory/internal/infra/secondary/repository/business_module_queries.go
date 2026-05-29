package repository

import (
	"context"

	"gorm.io/gorm"
)

func (r *Repository) IsBusinessModuleActive(ctx context.Context, businessID uint, moduleCode string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Table("integrations AS i").
		Joins("INNER JOIN integration_types AS it ON it.id = i.integration_type_id AND it.deleted_at IS NULL").
		Where("i.business_id = ?", businessID).
		Where("i.deleted_at IS NULL").
		Where("i.is_active = ?", true).
		Where("it.code = ?", moduleCode).
		Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
