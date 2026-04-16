package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// ============================================
// CARRIER RESOLUTION QUERY
// (Replicated locally — module isolation rule)
// Table consulted: integrations + integration_types + integration_categories
// ============================================

// GetActiveShippingCarrier finds the active shipping integration for a business.
// Returns nil (no error) if no carrier is configured.
func (r *Repository) GetActiveShippingCarrier(ctx context.Context, businessID uint) (*domain.CarrierInfo, error) {
	var result struct {
		IntegrationID     uint   `gorm:"column:integration_id"`
		IntegrationTypeID uint   `gorm:"column:integration_type_id"`
		ProviderCode      string `gorm:"column:provider_code"`
		BaseURL           string `gorm:"column:base_url"`
		IsTesting         bool   `gorm:"column:is_testing"`
		BaseURLTest       string `gorm:"column:base_url_test"`
	}

	err := r.db.Conn(ctx).
		Table("integrations i").
		Select("i.id AS integration_id, it.id AS integration_type_id, LOWER(it.code) AS provider_code, it.base_url AS base_url, i.is_testing AS is_testing, COALESCE(it.base_url_test, '') AS base_url_test").
		Joins("INNER JOIN integration_types it ON i.integration_type_id = it.id").
		Joins("INNER JOIN integration_categories ic ON it.category_id = ic.id").
		Where("ic.code = ?", "shipping").
		Where("i.business_id = ?", businessID).
		Where("i.is_active = ?", true).
		Where("i.deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	if result.IntegrationID == 0 {
		return nil, nil // no carrier configured
	}

	return &domain.CarrierInfo{
		IntegrationID:     result.IntegrationID,
		IntegrationTypeID: result.IntegrationTypeID,
		ProviderCode:      result.ProviderCode,
		BaseURL:           result.BaseURL,
		IsTesting:         result.IsTesting,
		BaseURLTest:       result.BaseURLTest,
	}, nil
}

// GetBusinessName retrieves the name of a business by its ID.
// Used to produce descriptive error messages when a carrier is not configured.
func (r *Repository) GetBusinessName(ctx context.Context, businessID uint) (string, error) {
	var result struct {
		Name string `gorm:"column:name"`
	}

	err := r.db.Conn(ctx).
		Table("business").
		Select("name").
		Where("id = ?", businessID).
		Limit(1).
		Scan(&result).Error

	if err != nil {
		return fmt.Sprintf("negocio #%d", businessID), nil
	}
	if result.Name == "" {
		return fmt.Sprintf("negocio #%d", businessID), nil
	}
	return result.Name, nil
}
