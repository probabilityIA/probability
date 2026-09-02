package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) FindSprintName(ctx context.Context, sprintID uint) (string, bool, error) {
	var result struct {
		Name string
	}
	err := r.db.Conn(ctx).Model(&models.Sprint{}).
		Select("name").
		Where("id = ?", sprintID).
		Limit(1).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", false, nil
		}
		return "", false, err
	}
	return result.Name, true, nil
}
