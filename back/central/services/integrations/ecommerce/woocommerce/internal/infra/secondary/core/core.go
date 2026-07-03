package core

import (
	"context"
	"fmt"
	"os"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases"
)

// WooCommerceCore implementa integrationcore.IIntegrationContract para WooCommerce.
// Embeds BaseIntegration para los métodos opcionales no soportados.
type WooCommerceCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IWooCommerceUseCase
}

// New crea el adaptador core de WooCommerce.
func New(useCase usecases.IWooCommerceUseCase) *WooCommerceCore {
	return &WooCommerceCore{useCase: useCase}
}

// TestConnection delega al use case la verificación de credenciales.
func (w *WooCommerceCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return w.useCase.TestConnection(ctx, config, credentials)
}

// SyncOrdersByIntegrationID sincroniza órdenes de WooCommerce (últimos 30 días).
func (w *WooCommerceCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return w.useCase.SyncOrders(ctx, integrationID)
}

// SyncOrdersByIntegrationIDWithParams sincroniza órdenes con parámetros personalizados.
func (w *WooCommerceCore) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	return w.useCase.SyncOrdersWithParams(ctx, integrationID, params)
}

func (w *WooCommerceCore) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	return w.useCase.UpdateInventory(ctx, integrationID, productExternalID, quantity)
}

// GetWebhookURL retorna la URL para los webhooks de WooCommerce.
func (w *WooCommerceCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/api/v1/woocommerce/webhook?integration_id=%d", baseURL, integrationID)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Probability crea estos webhooks automaticamente en tu tienda WooCommerce. " +
			"Recibiras las ordenes en tiempo real cuando se creen o actualicen.",
		Events: []string{
			"order.created",
			"order.updated",
		},
	}, nil
}

// ListWebhooks lista los webhooks de Probability registrados en la tienda.
func (w *WooCommerceCore) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	items, err := w.useCase.ListWebhooks(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	return result, nil
}

// DeleteWebhook elimina un webhook de la tienda WooCommerce.
func (w *WooCommerceCore) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	return w.useCase.DeleteWebhook(ctx, integrationID, webhookID)
}

// VerifyWebhooksByURL retorna los webhooks de Probability ya configurados en la tienda.
func (w *WooCommerceCore) VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]interface{}, error) {
	return w.ListWebhooks(ctx, integrationID)
}

// CreateWebhook crea los webhooks de ordenes en la tienda y retorna los configurados.
func (w *WooCommerceCore) CreateWebhook(ctx context.Context, integrationID string, baseURL string) (interface{}, error) {
	secret := os.Getenv("WOOCOMMERCE_WEBHOOK_SECRET")
	if err := w.useCase.CreateWebhooks(ctx, integrationID, baseURL, secret); err != nil {
		return nil, err
	}
	return w.ListWebhooks(ctx, integrationID)
}
