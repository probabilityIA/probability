package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
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

func (uc *jumpsellerUseCase) fetchIntegration(ctx context.Context, integrationID string) (*domain.Integration, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.ErrIntegrationNotFound
	}
	return integration, nil
}

func (uc *jumpsellerUseCase) buildCredential(ctx context.Context, integrationID string, integration *domain.Integration) (domain.Credential, error) {
	apiKey, err := uc.service.DecryptCredential(ctx, integrationID, "api_key")
	if err != nil {
		return domain.Credential{}, fmt.Errorf("decrypting api_key: %w", err)
	}
	if apiKey == "" {
		return domain.Credential{}, domain.ErrMissingAPIKey
	}

	apiSecret, err := uc.service.DecryptCredential(ctx, integrationID, "api_secret")
	if err != nil {
		return domain.Credential{}, fmt.Errorf("decrypting api_secret: %w", err)
	}
	if apiSecret == "" {
		return domain.Credential{}, domain.ErrMissingAPISecret
	}

	effectiveURL, err := resolveEffectiveBaseURL(integration)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Bool("is_testing", integration.IsTesting).
			Msg("El tipo de integracion Jumpseller no tiene la URL configurada en base de datos")
		return domain.Credential{}, err
	}

	return domain.Credential{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   effectiveURL,
	}, nil
}

func (uc *jumpsellerUseCase) resolveIntegration(ctx context.Context, integrationID string) (*domain.Integration, domain.Credential, error) {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return nil, domain.Credential{}, err
	}

	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return nil, domain.Credential{}, err
	}

	return integration, cred, nil
}

func (uc *jumpsellerUseCase) resolveIntegrationForBusiness(ctx context.Context, integrationID string, businessID uint) (*domain.Integration, domain.Credential, error) {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return nil, domain.Credential{}, err
	}

	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		uc.logger.Warn(ctx).
			Str("integration_id", integrationID).
			Uint("business_id", businessID).
			Msg("Intento de operar una integracion de Jumpseller que no pertenece al negocio")
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
