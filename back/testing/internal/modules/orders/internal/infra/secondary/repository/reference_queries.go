package repository

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/models"
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/domain/entities"
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

func (r *Repository) GetIntegrations(ctx context.Context, businessID uint) ([]entities.Integration, error) {
	var results []models.Integration

	err := r.db.Conn(ctx).
		Where("(business_id = ? OR business_id IS NULL) AND is_active = true", businessID).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	integrations := make([]entities.Integration, len(results))
	for i, ig := range results {
		integrations[i] = entities.Integration{
			ID:                ig.ID,
			Name:              ig.Name,
			Code:              ig.Code,
			Category:          ig.Category,
			IntegrationTypeID: ig.IntegrationTypeID,
		}
	}
	return integrations, nil
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
