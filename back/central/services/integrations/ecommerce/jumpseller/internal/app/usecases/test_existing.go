package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func (uc *jumpsellerUseCase) TestExistingConnection(ctx context.Context, integrationID string, businessID uint) (*domain.StoreInfo, error) {
	var cred domain.Credential
	var err error
	if businessID == 0 {
		_, cred, err = uc.resolveIntegration(ctx, integrationID)
	} else {
		_, cred, err = uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	}
	if err != nil {
		return nil, err
	}

	storeInfo, err := uc.client.GetStoreInfo(ctx, cred)
	if err != nil {
		uc.logger.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("Jumpseller test de conexion (existente) fallo")
		return nil, err
	}

	return storeInfo, nil
}
