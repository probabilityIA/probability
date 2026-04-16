package repository

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) GetClienteFinalRoleID(ctx context.Context) (uint, error) {
	var role models.Role
	err := r.db.Conn(ctx).
		Where("name = ? AND deleted_at IS NULL", "cliente_final").
		First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, domainerrors.ErrRoleNotFound
		}
		return 0, err
	}
	return role.ID, nil
}

func (r *Repository) GetRoleLevelByUserAndBusiness(ctx context.Context, userID, businessID uint) (int, error) {
	var result struct {
		Level int
	}
	err := r.db.Conn(ctx).
		Table("business_staff bs").
		Select("r.level").
		Joins("INNER JOIN roles r ON r.id = bs.role_id").
		Where("bs.user_id = ? AND bs.business_id = ? AND bs.deleted_at IS NULL AND r.deleted_at IS NULL", userID, businessID).
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Level, nil
}
