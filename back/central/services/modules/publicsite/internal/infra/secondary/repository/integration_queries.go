package repository

import (
	"context"

	"gorm.io/gorm"
)

// IsIntegrationActiveOrMissing checks if an integration of the given type is active for the business.
// If no integration exists, returns true for backward compatibility (rollout suave).
//
// Table consulted: integrations (managed by integrations/core module)
// Replicated locally to avoid sharing repositories between modules.
func (r *Repository) IsIntegrationActiveOrMissing(ctx context.Context, businessID uint, integrationTypeID uint) (bool, error) {
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
		return true, nil // No integration = allowed (backward compat)
	}
	if err != nil {
		return false, err
	}

	return result.IsActive, nil
}
