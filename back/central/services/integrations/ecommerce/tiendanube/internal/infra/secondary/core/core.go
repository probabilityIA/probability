package core

import (
	"context"
	"fmt"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
)

// TiendanubeCore implementa integrationcore.IIntegrationContract para Tiendanube.
// Embeds BaseIntegration para los metodos opcionales no soportados.
type TiendanubeCore struct {
	integrationcore.BaseIntegration
	useCase usecases.ITiendanubeUseCase
}

// New crea el adaptador core de Tiendanube.
func New(useCase usecases.ITiendanubeUseCase) *TiendanubeCore {
	return &TiendanubeCore{useCase: useCase}
}

// TestConnection delega al use case la verificacion de credenciales.
func (t *TiendanubeCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return t.useCase.TestConnection(ctx, config, credentials)
}

// GetWebhookURL retorna la URL para los webhooks de Tiendanube.
func (t *TiendanubeCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/integrations/tiendanube/webhook", baseURL)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Configura este webhook en el panel de Tiendanube > Apps > Webhooks. " +
			"Suscribete a los eventos de ordenes para recibir notificaciones en tiempo real.",
		Events: []string{
			"order/created",
			"order/updated",
			"order/cancelled",
		},
	}, nil
}
