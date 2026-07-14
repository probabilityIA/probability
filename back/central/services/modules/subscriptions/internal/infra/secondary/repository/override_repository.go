package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) CreateOverride(ctx context.Context, override *entities.BusinessModuleOverride) error {
	overrideDB := &models.BusinessModuleOverride{
		BusinessID:      override.BusinessID,
		ModuleCode:      override.ModuleCode,
		GrantedByUserID: override.GrantedByUserID,
		Notes:           override.Notes,
	}

	if err := r.db.Conn(ctx).Create(overrideDB).Error; err != nil {
		return err
	}

	override.ID = overrideDB.ID
	override.CreatedAt = overrideDB.CreatedAt
	return nil
}

func (r *Repository) DeleteOverride(ctx context.Context, businessID uint, moduleCode string) error {
	return r.db.Conn(ctx).
		Where("business_id = ? AND module_code = ?", businessID, moduleCode).
		Delete(&models.BusinessModuleOverride{}).Error
}

func (r *Repository) ListOverridesByBusiness(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error) {
	var overridesDB []models.BusinessModuleOverride
	if err := r.db.Conn(ctx).Where("business_id = ?", businessID).Find(&overridesDB).Error; err != nil {
		return nil, err
	}

	overrides := make([]entities.BusinessModuleOverride, len(overridesDB))
	for i, o := range overridesDB {
		overrides[i] = *overrideToEntity(&o)
	}
	return overrides, nil
}
