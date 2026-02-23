package core

import (
	"context"
	"fmt"

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

// GetWebhookURL retorna la URL para los webhooks de WooCommerce.
func (w *WooCommerceCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/integrations/woocommerce/webhook", baseURL)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Configura este webhook en WooCommerce > Ajustes > Avanzado > Webhooks. " +
			"Suscríbete a los eventos de órdenes para recibir notificaciones en tiempo real.",
		Events: []string{
			"order.created",
			"order.updated",
			"order.deleted",
			"order.restored",
		},
	}, nil
}
