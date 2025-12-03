package usecaseintegrations

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetIntegrationByType obtiene una integración por código de tipo y business_id, con credenciales desencriptadas
func (uc *IntegrationUseCase) GetIntegrationByType(ctx context.Context, integrationTypeCode string, businessID *uint) (*domain.IntegrationWithCredentials, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationByType")

	// Primero obtener el IntegrationType por código
	integrationType, err := uc.repo.GetIntegrationTypeByCode(ctx, integrationTypeCode)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Str("type_code", integrationTypeCode).
			Msg("Error al obtener tipo de integración por código")
		return nil, fmt.Errorf("%w '%s': %w", domain.ErrIntegrationTypeNotFound, integrationTypeCode, err)
	}

	// Obtener integración del repository usando el IntegrationTypeID
	integration, err := uc.repo.GetActiveIntegrationByIntegrationTypeID(ctx, integrationType.ID, businessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("integration_type_id", integrationType.ID).
			Str("type_code", integrationTypeCode).
			Msg("Error al obtener integración por tipo")
		return nil, err
	}

	// Desencriptar credenciales
	var decryptedCredentials domain.DecryptedCredentials
	if len(integration.Credentials) > 0 {
		decrypted, err := uc.encryption.DecryptCredentials(ctx, []byte(integration.Credentials))
		if err != nil {
			uc.log.Error(ctx).Err(err).
				Uint("id", integration.ID).
				Msg("Error al desencriptar credenciales")
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsDecrypt, err)
		}
		decryptedCredentials = decrypted
	}

	return &domain.IntegrationWithCredentials{
		Integration:          *integration,
		DecryptedCredentials: decryptedCredentials,
	}, nil
}
