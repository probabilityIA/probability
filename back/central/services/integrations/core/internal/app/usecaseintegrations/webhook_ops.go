package usecaseintegrations

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetWebhookURL obtiene la URL del webhook para una integración específica.
func (uc *IntegrationUseCase) GetWebhookURL(ctx context.Context, integrationID uint) (*domain.WebhookInfo, error) {
	ctx = log.WithFunctionCtx(ctx, "GetWebhookURL")

	provider, err := uc.getProviderForIntegration(ctx, fmt.Sprintf("%d", integrationID))
	if err != nil {
		return nil, err
	}
	baseURL, err := uc.getWebhookBaseURL()
	if err != nil {
		return nil, err
	}
	return provider.GetWebhookURL(ctx, baseURL, integrationID)
}

// ListWebhooks lista todos los webhooks de una integración.
func (uc *IntegrationUseCase) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	ctx = log.WithFunctionCtx(ctx, "ListWebhooks")

	provider, err := uc.getProviderForIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	return provider.ListWebhooks(ctx, integrationID)
}

// DeleteWebhook elimina un webhook de una integración.
func (uc *IntegrationUseCase) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	ctx = log.WithFunctionCtx(ctx, "DeleteWebhook")

	provider, err := uc.getProviderForIntegration(ctx, integrationID)
	if err != nil {
		return err
	}
	return provider.DeleteWebhook(ctx, integrationID, webhookID)
}

// VerifyWebhooksByURL verifica webhooks existentes que coincidan con nuestra URL.
func (uc *IntegrationUseCase) VerifyWebhooksByURL(ctx context.Context, integrationID string) ([]interface{}, error) {
	ctx = log.WithFunctionCtx(ctx, "VerifyWebhooksByURL")

	provider, err := uc.getProviderForIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	baseURL, err := uc.getWebhookBaseURL()
	if err != nil {
		return nil, err
	}
	return provider.VerifyWebhooksByURL(ctx, integrationID, baseURL)
}

// CreateWebhookForIntegration crea webhooks en la plataforma externa después de verificar y eliminar duplicados.
// Nombre diferente a IIntegrationContract.CreateWebhook para evitar colisión de firma.
func (uc *IntegrationUseCase) CreateWebhookForIntegration(ctx context.Context, integrationID string) (interface{}, error) {
	ctx = log.WithFunctionCtx(ctx, "CreateWebhookForIntegration")

	provider, err := uc.getProviderForIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	baseURL, err := uc.getWebhookBaseURL()
	if err != nil {
		return nil, err
	}
	return provider.CreateWebhook(ctx, integrationID, baseURL)
}

// getWebhookBaseURL obtiene la URL base para webhooks desde configuración.
func (uc *IntegrationUseCase) getWebhookBaseURL() (string, error) {
	baseURL := uc.config.Get("WEBHOOK_BASE_URL")
	if baseURL == "" {
		baseURL = uc.config.Get("URL_BASE_SWAGGER")
	}
	if baseURL == "" {
		return "", fmt.Errorf("WEBHOOK_BASE_URL o URL_BASE_SWAGGER no está configurada")
	}
	return baseURL, nil
}

// getProviderForIntegration obtiene el provider registrado para una integración por su ID.
func (uc *IntegrationUseCase) getProviderForIntegration(ctx context.Context, integrationID string) (domain.IIntegrationContract, error) {
	integration, err := uc.GetPublicIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener integración: %w", err)
	}
	provider, ok := uc.providerReg.Get(integration.IntegrationType)
	if !ok {
		return nil, fmt.Errorf("integración no registrada para tipo %d", integration.IntegrationType)
	}
	return provider, nil
}
