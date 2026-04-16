package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/entities"
)

func (r *Repository) GetProducts(ctx context.Context, businessID uint) ([]entities.Product, error) {
	var results []models.Product

	err := r.db.Conn(ctx).
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Limit(100).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	products := make([]entities.Product, len(results))
	for i, p := range results {
		products[i] = entities.Product{
			ID:       p.ID,
			Name:     p.Name,
			SKU:      p.SKU,
			Price:    p.Price,
			Currency: p.Currency,
		}
	}
	return products, nil
}

// GetIntegrations returns ecommerce (category_id=1) and platform (category_id=6) integrations
func (r *Repository) GetIntegrations(ctx context.Context, businessID uint) ([]entities.Integration, error) {
	type integrationRow struct {
		ID                  uint
		Name                string
		Code                string
		Category            string
		CategoryID          uint
		IntegrationTypeID   uint
		IntegrationTypeCode string
	}

	var results []integrationRow

	err := r.db.Conn(ctx).
		Table("integrations i").
		Select("i.id, i.name, i.code, i.category, it.category_id, i.integration_type_id, it.code as integration_type_code").
		Joins("JOIN integration_types it ON it.id = i.integration_type_id").
		Where("(i.business_id = ? OR i.business_id IS NULL)", businessID).
		Where("i.is_active = true").
		Where("i.deleted_at IS NULL").
		Where("it.category_id IN ?", []uint{1, 6}). // ecommerce + platform
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	integrations := make([]entities.Integration, len(results))
	for i, ig := range results {
		integrations[i] = entities.Integration{
			ID:                  ig.ID,
			Name:                ig.Name,
			Code:                ig.Code,
			Category:            ig.Category,
			CategoryID:          ig.CategoryID,
			IntegrationTypeID:   ig.IntegrationTypeID,
			IntegrationTypeCode: ig.IntegrationTypeCode,
		}
	}
	return integrations, nil
}

// GetIntegrationTypeCode returns the integration_type code for a given integration ID
func (r *Repository) GetIntegrationTypeCode(ctx context.Context, integrationID uint) (string, error) {
	var code string
	err := r.db.Conn(ctx).
		Table("integrations i").
		Select("it.code").
		Joins("JOIN integration_types it ON it.id = i.integration_type_id").
		Where("i.id = ?", integrationID).
		Scan(&code).Error
	if err != nil {
		return "", fmt.Errorf("failed to get integration type code: %w", err)
	}
	return code, nil
}

// GetIntegrationCategoryID returns the category_id of the integration type for a given integration ID
func (r *Repository) GetIntegrationCategoryID(ctx context.Context, integrationID uint) (uint, error) {
	var categoryID uint
	err := r.db.Conn(ctx).
		Table("integrations i").
		Select("it.category_id").
		Joins("JOIN integration_types it ON it.id = i.integration_type_id").
		Where("i.id = ?", integrationID).
		Scan(&categoryID).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get integration category: %w", err)
	}
	return categoryID, nil
}

func (r *Repository) GetPaymentMethods(ctx context.Context) ([]entities.PaymentMethod, error) {
	var results []models.PaymentMethod

	err := r.db.Conn(ctx).
		Where("is_active = true").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	methods := make([]entities.PaymentMethod, len(results))
	for i, pm := range results {
		methods[i] = entities.PaymentMethod{
			ID:   pm.ID,
			Code: pm.Code,
			Name: pm.Name,
		}
	}
	return methods, nil
}

func (r *Repository) GetOrderStatuses(ctx context.Context) ([]entities.OrderStatus, error) {
	var results []models.OrderStatus

	err := r.db.Conn(ctx).
		Where("is_active = true").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	statuses := make([]entities.OrderStatus, len(results))
	for i, os := range results {
		statuses[i] = entities.OrderStatus{
			ID:   os.ID,
			Code: os.Code,
			Name: os.Name,
		}
	}
	return statuses, nil
}
