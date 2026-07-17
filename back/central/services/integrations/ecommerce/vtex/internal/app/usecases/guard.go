package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

func (uc *vtexUseCase) AssertIntegrationOwned(ctx context.Context, integrationID string, businessID uint) error {
	_, err := uc.integrationForBusiness(ctx, integrationID, businessID)
	return err
}

func (uc *vtexUseCase) integrationForBusiness(ctx context.Context, integrationID string, businessID uint) (*domain.Integration, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.ErrIntegrationNotFound
	}

	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		uc.logger.Warn(ctx).
			Str("integration_id", integrationID).
			Uint("business_id", businessID).
			Msg("Intento de operar sobre una integracion VTEX de otro negocio")
		return nil, domain.ErrIntegrationNotOwned
	}

	return integration, nil
}
