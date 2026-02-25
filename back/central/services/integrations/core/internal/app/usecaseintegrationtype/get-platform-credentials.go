package usecaseintegrationtype

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetPlatformCredentials decrypts and returns the platform credentials for the given integration type.
// Intended for admin use only — never exposed publicly.
func (uc *integrationTypeUseCase) GetPlatformCredentials(ctx context.Context, id uint) (map[string]interface{}, error) {
	ctx = log.WithFunctionCtx(ctx, "GetPlatformCredentials")

	integrationType, err := uc.repo.GetIntegrationTypeByID(ctx, id)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener tipo de integración")
		return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTypeNotFound, err)
	}

	if len(integrationType.PlatformCredentialsEncrypted) == 0 {
		return map[string]interface{}{}, nil
	}

	creds, err := uc.encryption.DecryptCredentials(ctx, integrationType.PlatformCredentialsEncrypted)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al desencriptar credenciales de plataforma")
		return nil, fmt.Errorf("error al desencriptar credenciales: %w", err)
	}

	return creds, nil
}
