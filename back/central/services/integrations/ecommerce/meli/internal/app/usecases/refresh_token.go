package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) EnsureValidToken(ctx context.Context, integrationID string) (string, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return "", fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return "", domain.ErrIntegrationNotFound
	}

	needsRefresh := true
	if expiresAtStr, ok := integration.Config["token_expires_at"].(string); ok && expiresAtStr != "" {
		expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
		if err == nil {
			needsRefresh = time.Now().After(expiresAt.Add(-5 * time.Minute))
		}
	}

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

	accessToken, err := uc.service.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return "", fmt.Errorf("decrypting access_token: %w", err)
	}
	if accessToken == "" {
		return "", domain.ErrMissingAccessToken
	}

	return accessToken, nil
}

func (uc *meliUseCase) refreshAccessToken(ctx context.Context, integrationID string, integration *domain.Integration) (string, error) {
	appID, err := extractStringFromConfig(integration.Config, "app_id")
	if err != nil {
		return "", domain.ErrMissingAppID
	}

	clientSecret, err := uc.service.DecryptCredential(ctx, integrationID, "client_secret")
	if err != nil {
		return "", fmt.Errorf("decrypting client_secret: %w", err)
	}
	if clientSecret == "" {
		return "", domain.ErrMissingClientSecret
	}

	refreshToken, err := uc.service.DecryptCredential(ctx, integrationID, "refresh_token")
	if err != nil {
		return "", fmt.Errorf("decrypting refresh_token: %w", err)
	}
	if refreshToken == "" {
		return "", domain.ErrMissingRefreshToken
	}

	tokenResp, err := uc.client.RefreshToken(ctx, appID, clientSecret, refreshToken)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Msg("Token refresh API call failed")
		return "", fmt.Errorf("%w: %v", domain.ErrTokenRefreshFailed, err)
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	newRefreshToken := refreshToken
	if tokenResp.RefreshToken != "" {
		newRefreshToken = tokenResp.RefreshToken
	}

	newCredentials := map[string]interface{}{
		"access_token":  tokenResp.AccessToken,
		"refresh_token": newRefreshToken,
		"client_secret": clientSecret,
	}
	if err := uc.service.UpdateIntegrationCredentials(ctx, integrationID, newCredentials); err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Msg("Failed to persist refreshed MeLi credentials")
	}

	newConfig := make(map[string]interface{})
	for k, v := range integration.Config {
		newConfig[k] = v
	}
	newConfig["token_expires_at"] = expiresAt.Format(time.RFC3339)
	if tokenResp.UserID > 0 {
		newConfig["seller_id"] = tokenResp.UserID
	}

	if err := uc.service.UpdateIntegrationConfig(ctx, integrationID, newConfig); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("Failed to update integration config after token refresh")
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("expires_at", expiresAt.Format(time.RFC3339)).
		Msg("MeLi token refreshed successfully")

	return tokenResp.AccessToken, nil
}

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
