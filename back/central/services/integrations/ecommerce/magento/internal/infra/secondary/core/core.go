package core

import (
	"context"
	"fmt"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/magento/internal/app/usecases"
)

// MagentoCore implementa integrationcore.IIntegrationContract para Magento.
// Embeds BaseIntegration para los métodos opcionales no soportados.
type MagentoCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IMagentoUseCase
}

// New crea el adaptador core de Magento.
func New(useCase usecases.IMagentoUseCase) *MagentoCore {
	return &MagentoCore{useCase: useCase}
}

// TestConnection delega al use case la verificación de credenciales.
func (m *MagentoCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return m.useCase.TestConnection(ctx, config, credentials)
}

// GetWebhookURL retorna la URL para los webhooks de Magento.
func (m *MagentoCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/integrations/magento/webhook", baseURL)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Configura este webhook en Magento > Sistema > Integraciones o mediante una extensión de webhooks. " +
			"Suscríbete a los eventos de órdenes para recibir notificaciones en tiempo real.",
		Events: []string{
			"sales_order_save_after",
			"sales_order_place_after",
		},
	}, nil
}
