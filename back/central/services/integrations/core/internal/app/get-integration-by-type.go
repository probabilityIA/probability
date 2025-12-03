package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetIntegrationByType obtiene una integración por tipo y business_id, con credenciales desencriptadas
func (uc *integrationUseCase) GetIntegrationByType(ctx context.Context, integrationType string, businessID *uint) (*domain.IntegrationWithCredentials, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationByType")

	// Obtener integración del repository (credenciales encriptadas)
	integration, err := uc.repo.GetActiveByType(ctx, integrationType, businessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Str("type", integrationType).
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
			return nil, fmt.Errorf("error al desencriptar credenciales: %w", err)
		}
		decryptedCredentials = decrypted
	}

	return &domain.IntegrationWithCredentials{
		Integration:          *integration,
		DecryptedCredentials: decryptedCredentials,
	}, nil
}
