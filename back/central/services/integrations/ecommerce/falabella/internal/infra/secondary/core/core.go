package core

import (
	"context"
	"fmt"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/falabella/internal/app/usecases"
)

// FalabellaCore implementa integrationcore.IIntegrationContract para Falabella.
// Embeds BaseIntegration para los métodos opcionales no soportados.
type FalabellaCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IFalabellaUseCase
}

// New crea el adaptador core de Falabella.
func New(useCase usecases.IFalabellaUseCase) *FalabellaCore {
	return &FalabellaCore{useCase: useCase}
}

// TestConnection delega al use case la verificación de credenciales.
func (f *FalabellaCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return f.useCase.TestConnection(ctx, config, credentials)
}

// GetWebhookURL retorna la URL para los webhooks de Falabella Seller Center.
func (f *FalabellaCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/integrations/falabella/webhook", baseURL)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Configura este webhook en Falabella Seller Center para recibir " +
			"notificaciones de órdenes en tiempo real.",
		Events: []string{
			"order.created",
			"order.updated",
		},
	}, nil
}
