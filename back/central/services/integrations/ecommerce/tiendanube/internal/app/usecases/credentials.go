package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

const (
	configStoreID   = "store_id"
	configStoreURL  = "store_url"
	configUserAgent = "user_agent"
)

func extractString(m map[string]interface{}, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("missing field: %s", key)
	}
	s, ok := v.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("field %s must be a non-empty string", key)
	}
	return strings.TrimSpace(s), nil
}

func optionalString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}

func resolveStoreID(config map[string]interface{}, fallback string) string {
	if id := optionalString(config, configStoreID); id != "" {
		return id
	}
	if url := optionalString(config, configStoreURL); url != "" {
		if id := storeIDFromURL(url); id != "" {
			return id
		}
	}
	return strings.TrimSpace(fallback)
}

func storeIDFromURL(raw string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(raw), "/")
	if trimmed == "" {
		return ""
	}
	parts := strings.Split(trimmed, "/")
	last := parts[len(parts)-1]
	for _, r := range last {
		if r < '0' || r > '9' {
			return ""
		}
	}
	return last
}

func resolveEffectiveBaseURL(integration *domain.Integration) (string, error) {
	if integration.IsTesting {
		if strings.TrimSpace(integration.BaseURLTest) == "" {
			return "", domain.ErrMissingBaseURLTest
		}
		return integration.BaseURLTest, nil
	}
	if strings.TrimSpace(integration.BaseURL) == "" {
		return "", domain.ErrMissingBaseURL
	}
	return integration.BaseURL, nil
}

func (uc *tiendanubeUseCase) fetchIntegration(ctx context.Context, integrationID string) (*domain.Integration, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.ErrIntegrationNotFound
	}
	return integration, nil
}

func (uc *tiendanubeUseCase) buildCredential(ctx context.Context, integrationID string, integration *domain.Integration) (domain.Credential, error) {
	accessToken, err := uc.service.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return domain.Credential{}, fmt.Errorf("decrypting access_token: %w", err)
	}
	if strings.TrimSpace(accessToken) == "" {
		return domain.Credential{}, domain.ErrMissingAccessToken
	}

	storeID := resolveStoreID(integration.Config, integration.StoreID)
	if storeID == "" {
		return domain.Credential{}, domain.ErrMissingStoreID
	}

	baseURL, err := resolveEffectiveBaseURL(integration)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Bool("is_testing", integration.IsTesting).
			Msg("El tipo de integracion Tiendanube no tiene la URL configurada en base de datos")
		return domain.Credential{}, err
	}

	return domain.Credential{
		AccessToken: accessToken,
		StoreID:     storeID,
		BaseURL:     baseURL,
		UserAgent:   optionalString(integration.Config, configUserAgent),
	}, nil
}

func (uc *tiendanubeUseCase) resolveIntegrationForBusiness(ctx context.Context, integrationID string, businessID uint) (*domain.Integration, domain.Credential, error) {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return nil, domain.Credential{}, err
	}

	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		uc.logger.Warn(ctx).
			Str("integration_id", integrationID).
			Uint("business_id", businessID).
			Msg("Intento de operar una integracion de Tiendanube que no pertenece al negocio")
		return nil, domain.Credential{}, domain.ErrIntegrationNotFound
	}

	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return nil, domain.Credential{}, err
	}

	return integration, cred, nil
}
