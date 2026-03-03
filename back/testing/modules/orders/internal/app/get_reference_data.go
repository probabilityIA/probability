package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/entities"
)

func (uc *useCase) GetReferenceData(ctx context.Context, businessID uint) (*entities.ReferenceData, error) {
	products, err := uc.repo.GetProducts(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	integrations, err := uc.repo.GetIntegrations(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integrations: %w", err)
	}

	paymentMethods, err := uc.repo.GetPaymentMethods(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	orderStatuses, err := uc.repo.GetOrderStatuses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get order statuses: %w", err)
	}

	return &entities.ReferenceData{
		Products:       products,
		Integrations:   integrations,
		PaymentMethods: paymentMethods,
		OrderStatuses:  orderStatuses,
	}, nil
}
