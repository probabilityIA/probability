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

	// 1. Obtener metadata de integración (con cache)
	integration, err := uc.GetIntegrationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationNotFound, err)
	}

	// 2. ✅ NUEVO - Intentar leer credentials de cache
	cachedCreds, err := uc.cache.GetCredentials(ctx, id)
	if err == nil {
		uc.log.Debug(ctx).Uint("id", id).Msg("✅ Cache hit - credentials")
		return &domain.IntegrationWithCredentials{
			Integration:          *integration,
			DecryptedCredentials: cachedCreds.Credentials,
		}, nil
	}

	// 3. Cache miss - Desencriptar credenciales de DB
	uc.log.Debug(ctx).Uint("id", id).Msg("⚠️ Cache miss credentials - decrypting")
	var decryptedCredentials domain.DecryptedCredentials
	if len(integration.Credentials) > 0 {
		// Procesamos credenciales sin loguearlas (seguridad)

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

		// Credenciales desencriptadas exitosamente (no logueamos keys por seguridad)
		decryptedCredentials = decrypted

		// ✅ NUEVO - Cachear credentials para próxima vez
		cachedCreds := &domain.CachedCredentials{
			IntegrationID: id,
			Credentials:   decryptedCredentials,
		}
		if err := uc.cache.SetCredentials(ctx, cachedCreds); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Failed to cache credentials")
		}
	}

	return &domain.IntegrationWithCredentials{
		Integration:          *integration,
		DecryptedCredentials: decryptedCredentials,
	}, nil
}
