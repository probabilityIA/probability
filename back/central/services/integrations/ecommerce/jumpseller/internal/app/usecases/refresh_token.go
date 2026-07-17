package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func (uc *jumpsellerUseCase) EnsureValidToken(ctx context.Context, integrationID string, integration *domain.Integration) (string, error) {
	needsRefresh := true
	if expiresAtStr, ok := integration.Config[domain.ConfigTokenExpiresAt].(string); ok && expiresAtStr != "" {
		if expiresAt, err := time.Parse(time.RFC3339, expiresAtStr); err == nil {
			needsRefresh = time.Now().After(expiresAt.Add(-5 * time.Minute))
		}
	}

	if needsRefresh {
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Msg("Token de Jumpseller expirado o por expirar, renovando")
		return uc.refreshAccessToken(ctx, integrationID)
	}

	accessToken, err := uc.service.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return "", fmt.Errorf("decrypting access_token: %w", err)
	}
	if accessToken == "" {
		return "", domain.ErrMissingAccessToken
	}
	return accessToken, nil
}

func (uc *jumpsellerUseCase) refreshAccessToken(ctx context.Context, integrationID string) (string, error) {
	clientID, err := uc.service.GetPlatformCredential(ctx, integrationID, "client_id")
	if err != nil || clientID == "" {
		return "", domain.ErrMissingClientID
	}

	clientSecret, err := uc.service.GetPlatformCredential(ctx, integrationID, "client_secret")
	if err != nil || clientSecret == "" {
		return "", domain.ErrMissingClientSecret
	}

	refreshToken, err := uc.service.DecryptCredential(ctx, integrationID, "refresh_token")
	if err != nil {
		return "", fmt.Errorf("decrypting refresh_token: %w", err)
	}
	if refreshToken == "" {
		return "", domain.ErrMissingRefreshToken
	}

	tokenResp, err := uc.client.RefreshToken(ctx, domain.OAuthTokenURL, clientID, clientSecret, refreshToken)
	if err != nil {
		uc.logger.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("Fallo la renovacion del token de Jumpseller")
		return "", err
	}

	newRefreshToken := refreshToken
	if tokenResp.RefreshToken != "" {
		newRefreshToken = tokenResp.RefreshToken
	}

	if err := uc.service.UpdateIntegrationCredentials(ctx, integrationID, map[string]interface{}{
		"access_token":  tokenResp.AccessToken,
		"refresh_token": newRefreshToken,
	}); err != nil {
		uc.logger.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("No se pudieron persistir las credenciales renovadas de Jumpseller")
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	if err := uc.service.UpdateIntegrationConfig(ctx, integrationID, map[string]interface{}{
		domain.ConfigTokenExpiresAt: expiresAt.Format(time.RFC3339),
	}); err != nil {
		uc.logger.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("No se pudo actualizar token_expires_at de Jumpseller")
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("expires_at", expiresAt.Format(time.RFC3339)).
		Msg("Token de Jumpseller renovado")

	return tokenResp.AccessToken, nil
}
