package core

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases"
)

type ShopifyCore struct {
	core.BaseIntegration
	useCase usecases.IShopifyUseCase
}

func New(useCase usecases.IShopifyUseCase) *ShopifyCore {
	return &ShopifyCore{useCase: useCase}
}

func (s *ShopifyCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return s.useCase.TestConnection(ctx, config, credentials)
}

func (s *ShopifyCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return s.useCase.SyncOrders(ctx, integrationID)
}

// GetWebhookURL construye la URL del webhook para Shopify
func (s *ShopifyCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*core.WebhookInfo, error) {
	// Construir la URL del webhook
	// El webhook se recibe en: /api/v1/integrations/shopify/webhook
	webhookURL := fmt.Sprintf("%s/integrations/shopify/webhook", baseURL)

	return &core.WebhookInfo{
		URL:         webhookURL,
		Method:      "POST",
		Description: "URL para configurar en Shopify Admin > Settings > Notifications > Webhooks. Configure este webhook para recibir eventos de órdenes en tiempo real.",
		Events: []string{
			"orders/create",
			"orders/updated",
			"orders/paid",
			"orders/cancelled",
			"orders/fulfilled",
			"orders/partially_fulfilled",
		},
	}, nil
}

// ListWebhooks lista todos los webhooks de una integración de Shopify
func (s *ShopifyCore) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	webhooks, err := s.useCase.ListWebhooks(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	// Convertir a []interface{} para la interfaz
	result := make([]interface{}, len(webhooks))
	for i, wh := range webhooks {
		result[i] = wh
	}

	return result, nil
}

// DeleteWebhook elimina un webhook de Shopify
func (s *ShopifyCore) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	return s.useCase.DeleteWebhook(ctx, integrationID, webhookID)
}

// VerifyWebhooksByURL verifica webhooks existentes que coincidan con nuestra URL
func (s *ShopifyCore) VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]interface{}, error) {
	webhooks, err := s.useCase.VerifyWebhooksByURL(ctx, integrationID, baseURL)
	if err != nil {
		return nil, err
	}

	// Convertir a []interface{} para la interfaz
	result := make([]interface{}, len(webhooks))
	for i, wh := range webhooks {
		result[i] = wh
	}

	return result, nil
}

// CreateWebhook crea webhooks en Shopify después de verificar y eliminar duplicados
func (s *ShopifyCore) CreateWebhook(ctx context.Context, integrationID string, baseURL string) (interface{}, error) {
	result, err := s.useCase.CreateWebhook(ctx, integrationID, baseURL)
	if err != nil {
		return nil, err
	}

	return result, nil
}
