package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

// EnsureValidToken verifica si el token actual está vigente, y si no, lo renueva.
// Usa un margen de 5 minutos para renovar antes de la expiración real.
// Retorna un access_token listo para usar.
func (uc *meliUseCase) EnsureValidToken(ctx context.Context, integrationID string) (string, error) {
	// 1. Obtener integración para revisar config["token_expires_at"]
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return "", fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return "", domain.ErrIntegrationNotFound
	}

	// 2. Verificar si el token necesita renovación
	needsRefresh := true
	if expiresAtStr, ok := integration.Config["token_expires_at"].(string); ok && expiresAtStr != "" {
		expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
		if err == nil {
			// Renovar si faltan menos de 5 minutos para expirar
			needsRefresh = time.Now().After(expiresAt.Add(-5 * time.Minute))
		}
	}

	// 3. Si necesita refresh, renovar
	if needsRefresh {
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Msg("Token expired or expiring soon, refreshing...")

		newToken, err := uc.refreshAccessToken(ctx, integrationID, integration)
		if err != nil {
			return "", err
		}
		return newToken, nil
	}

	// 4. Token vigente — desencriptar y retornar
	accessToken, err := uc.service.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return "", fmt.Errorf("decrypting access_token: %w", err)
	}
	if accessToken == "" {
		return "", domain.ErrMissingAccessToken
	}

	return accessToken, nil
}

// refreshAccessToken renueva el token usando el refresh_token y guarda los nuevos valores.
func (uc *meliUseCase) refreshAccessToken(ctx context.Context, integrationID string, integration *domain.Integration) (string, error) {
	// 1. Obtener app_id del config
	appID, err := extractStringFromConfig(integration.Config, "app_id")
	if err != nil {
		return "", domain.ErrMissingAppID
	}

	// 2. Desencriptar client_secret
	clientSecret, err := uc.service.DecryptCredential(ctx, integrationID, "client_secret")
	if err != nil {
		return "", fmt.Errorf("decrypting client_secret: %w", err)
	}
	if clientSecret == "" {
		return "", domain.ErrMissingClientSecret
	}

	// 3. Desencriptar refresh_token
	refreshToken, err := uc.service.DecryptCredential(ctx, integrationID, "refresh_token")
	if err != nil {
		return "", fmt.Errorf("decrypting refresh_token: %w", err)
	}
	if refreshToken == "" {
		return "", domain.ErrMissingRefreshToken
	}

	// 4. Llamar a la API para refrescar
	tokenResp, err := uc.client.RefreshToken(ctx, appID, clientSecret, refreshToken)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Msg("Token refresh API call failed")
		return "", fmt.Errorf("%w: %v", domain.ErrTokenRefreshFailed, err)
	}

	// 5. Calcular expiración
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// 6. Actualizar config con el nuevo token_expires_at
	newConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		newConfig[k] = v
	}
	newConfig["token_expires_at"] = expiresAt.Format(time.RFC3339)
	// Almacenar seller_id si viene en la respuesta del token
	if tokenResp.UserID > 0 {
		newConfig["seller_id"] = tokenResp.UserID
	}

	if err := uc.service.UpdateIntegrationConfig(ctx, integrationID, newConfig); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("Failed to update integration config after token refresh")
		// No retornar error — el token es válido aunque no se haya guardado el config
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("expires_at", expiresAt.Format(time.RFC3339)).
		Msg("MeLi token refreshed successfully")

	return tokenResp.AccessToken, nil
}

// extractStringFromConfig extrae un campo string de un mapa de config.
func extractStringFromConfig(config map[string]interface{}, key string) (string, error) {
	v, ok := config[key]
	if !ok {
		return "", fmt.Errorf("missing config field: %s", key)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("config field %s must be a non-empty string", key)
	}
	return s, nil
}
