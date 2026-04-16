package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) ListEcommerceIntegrationTypes(ctx context.Context, businessID uint) ([]entities.IntegrationTypeInfo, error) {
	return uc.repo.ListEcommerceIntegrationTypes(ctx, businessID)
}
