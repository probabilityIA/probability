package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
)

func resolveEffectiveBaseURL(integration *domain.Integration) (string, error) {
	if integration.IsTesting {
		if integration.BaseURLTest == "" {
			return "", domain.ErrMissingBaseURLTest
		}
		return integration.BaseURLTest, nil
	}
	if integration.BaseURL == "" {
		return "", domain.ErrMissingBaseURL
	}
	return integration.BaseURL, nil
}

func testConnectionBaseURL(config map[string]interface{}) (string, error) {
	if isTesting, _ := config["is_testing"].(bool); isTesting {
		url, err := extractString(config, "base_url_test")
		if err != nil {
			return "", domain.ErrMissingBaseURLTest
		}
		return url, nil
	}
	url, err := extractString(config, "base_url")
	if err != nil {
		return "", domain.ErrMissingBaseURL
	}
	return url, nil
}

func (uc *tiktokUseCase) fetchIntegration(ctx context.Context, integrationID string) (*domain.Integration, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.ErrIntegrationNotFound
	}
	return integration, nil
}

func (uc *tiktokUseCase) buildCredential(ctx context.Context, integrationID string, integration *domain.Integration) (domain.Credential, error) {
	clientKey, err := uc.service.DecryptCredential(ctx, integrationID, "client_key")
	if err != nil {
		return domain.Credential{}, fmt.Errorf("decrypting client_key: %w", err)
	}
	if clientKey == "" {
		return domain.Credential{}, domain.ErrMissingClientKey
	}

	clientSecret, err := uc.service.DecryptCredential(ctx, integrationID, "client_secret")
	if err != nil {
		return domain.Credential{}, fmt.Errorf("decrypting client_secret: %w", err)
	}
	if clientSecret == "" {
		return domain.Credential{}, domain.ErrMissingClientSecret
	}

	effectiveURL, err := resolveEffectiveBaseURL(integration)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Bool("is_testing", integration.IsTesting).
			Msg("El tipo de integracion TikTok no tiene la URL configurada en base de datos")
		return domain.Credential{}, err
	}

	token, err := uc.client.GetClientToken(ctx, effectiveURL, clientKey, clientSecret)
	if err != nil {
		return domain.Credential{}, fmt.Errorf("getting client token: %w", err)
	}

	return domain.Credential{
		AccessToken: token,
		BaseURL:     effectiveURL,
	}, nil
}

func (uc *tiktokUseCase) resolveIntegrationForBusiness(ctx context.Context, integrationID string, businessID uint) (*domain.Integration, domain.Credential, error) {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return nil, domain.Credential{}, err
	}

	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		uc.logger.Warn(ctx).
			Str("integration_id", integrationID).
			Uint("business_id", businessID).
			Msg("Intento de operar una integracion de TikTok que no pertenece al negocio")
		return nil, domain.Credential{}, domain.ErrIntegrationNotFound
	}

	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return nil, domain.Credential{}, err
	}

	return integration, cred, nil
}

func extractString(m map[string]interface{}, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("missing field: %s", key)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("field %s must be a non-empty string", key)
	}
	return s, nil
}
