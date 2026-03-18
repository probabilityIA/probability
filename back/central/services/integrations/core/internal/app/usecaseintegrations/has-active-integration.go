package usecaseintegrations

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/log"
)

// HasActiveIntegration verifica si existe una integración activa por integration_type_id y business_id
func (uc *IntegrationUseCase) HasActiveIntegration(ctx context.Context, integrationTypeID uint, businessID *uint) (bool, error) {
	ctx = log.WithFunctionCtx(ctx, "HasActiveIntegration")

	exists, err := uc.repo.ExistsActiveIntegrationByTypeID(ctx, integrationTypeID, businessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("integration_type_id", integrationTypeID).
			Msg("Error verificando existencia de integración activa")
		return false, err
	}

	return exists, nil
}

// HasActiveIntegrationByCode verifica si existe una integración activa resolviendo primero el código de tipo
func (uc *IntegrationUseCase) HasActiveIntegrationByCode(ctx context.Context, integrationTypeCode string, businessID *uint) (bool, error) {
	ctx = log.WithFunctionCtx(ctx, "HasActiveIntegrationByCode")

	integrationType, err := uc.repo.GetIntegrationTypeByCode(ctx, integrationTypeCode)
	if err != nil {
		return false, nil
	}

	return uc.repo.ExistsActiveIntegrationByTypeID(ctx, integrationType.ID, businessID)
}
