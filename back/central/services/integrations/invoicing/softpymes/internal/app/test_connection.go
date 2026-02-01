package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/errors"
)

// TestProviderConnection prueba la conexión con el proveedor Softpymes
func (uc *useCase) TestProviderConnection(ctx context.Context, id uint) error {
	uc.log.Info(ctx).Uint("provider_id", id).Msg("Testing Softpymes provider connection")

	// 1. Obtener proveedor
	provider, err := uc.providerRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrProviderNotFound
	}

	// 2. Extraer credenciales y configuración
	apiKey, okKey := provider.Credentials["api_key"].(string)
	apiSecret, okSecret := provider.Credentials["api_secret"].(string)
	referer, okReferer := provider.Config["referer"].(string)

	if !okKey || apiKey == "" {
		return errors.ErrAPIKeyRequired
	}

	if !okSecret || apiSecret == "" {
		return errors.ErrAPISecretRequired
	}

	if !okReferer || referer == "" {
		uc.log.Error(ctx).Msg("Referer is missing in provider config")
		return errors.ErrAuthenticationFailed // TODO: crear error específico para referer
	}

	// 3. Probar autenticación con Softpymes
	if err := uc.softpymesClient.TestAuthentication(ctx, apiKey, apiSecret, referer); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Softpymes connection test failed")
		return errors.ErrAuthenticationFailed
	}

	uc.log.Info(ctx).Uint("provider_id", id).Msg("Softpymes connection test successful")
	return nil
}
