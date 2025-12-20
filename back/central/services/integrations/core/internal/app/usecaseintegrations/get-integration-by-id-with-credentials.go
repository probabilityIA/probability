package usecaseintegrations

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetIntegrationByIDWithCredentials obtiene una integración por su ID con credenciales desencriptadas
// Este método es para uso en edición, solo debe ser accesible por super admins
func (uc *IntegrationUseCase) GetIntegrationByIDWithCredentials(ctx context.Context, id uint) (*domain.IntegrationWithCredentials, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationByIDWithCredentials")

	// Obtener integración
	integration, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener integración")
		return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationNotFound, err)
	}

	// Desencriptar credenciales
	var decryptedCredentials domain.DecryptedCredentials
	if len(integration.Credentials) > 0 {
		// DEBUG: Log raw credentials from DB
		uc.log.Debug(ctx).
			Uint("id", integration.ID).
			Str("raw_credentials", string(integration.Credentials)).
			Msg("Raw credentials from database")

		// Las credenciales están codificadas en base64 dentro de un JSON
		encryptedBytes, err := decodeEncryptedCredentials([]byte(integration.Credentials))
		if err != nil {
			uc.log.Error(ctx).Err(err).
				Uint("id", integration.ID).
				Msg("Error al decodificar credenciales desde base64")
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsDecrypt, err)
		}
		decrypted, err := uc.encryption.DecryptCredentials(ctx, encryptedBytes)
		if err != nil {
			uc.log.Error(ctx).Err(err).
				Uint("id", integration.ID).
				Msg("Error al desencriptar credenciales")
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsDecrypt, err)
		}

		// DEBUG: Log decrypted credential keys
		keys := make([]string, 0, len(decrypted))
		for k := range decrypted {
			keys = append(keys, k)
		}
		uc.log.Debug(ctx).
			Uint("id", integration.ID).
			Strs("decrypted_keys", keys).
			Msg("Decrypted credential keys")

		decryptedCredentials = decrypted
	}

	return &domain.IntegrationWithCredentials{
		Integration:          *integration,
		DecryptedCredentials: decryptedCredentials,
	}, nil
}
