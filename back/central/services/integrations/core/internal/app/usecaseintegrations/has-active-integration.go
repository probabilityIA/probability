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

// HasActiveIntegrationByCode verifica si existe una integración activa resolviendo primero el código de tipo.
// Busca case-insensitive y maneja variantes comunes (whatsapp, whatsap, whastap).
func (uc *IntegrationUseCase) HasActiveIntegrationByCode(ctx context.Context, integrationTypeCode string, businessID *uint) (bool, error) {
	ctx = log.WithFunctionCtx(ctx, "HasActiveIntegrationByCode")

	// Buscar por código exacto primero
	integrationType, err := uc.repo.GetIntegrationTypeByCode(ctx, integrationTypeCode)
	if err != nil {
		// Fallback: buscar case-insensitive
		integrationType, err = uc.repo.GetIntegrationTypeByCodeInsensitive(ctx, integrationTypeCode)
		if err != nil {
			uc.log.Warn(ctx).
				Str("code", integrationTypeCode).
				Msg("Tipo de integración no encontrado por código")
			return false, nil
		}
	}

	return uc.repo.ExistsActiveIntegrationByTypeID(ctx, integrationType.ID, businessID)
}
