package core

import (
	"context"
	"fmt"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/exito/internal/app/usecases"
)

// ExitoCore implementa integrationcore.IIntegrationContract para Exito.
// Embeds BaseIntegration para los metodos opcionales no soportados.
type ExitoCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IExitoUseCase
}

// New crea el adaptador core de Exito.
func New(useCase usecases.IExitoUseCase) *ExitoCore {
	return &ExitoCore{useCase: useCase}
}

// TestConnection delega al use case la verificacion de credenciales.
func (e *ExitoCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return e.useCase.TestConnection(ctx, config, credentials)
}

// GetWebhookURL retorna la URL para los webhooks de Exito.
func (e *ExitoCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/integrations/exito/webhook", baseURL)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Configura este webhook en el panel de seller de Exito. " +
			"Suscribete a los eventos de ordenes para recibir notificaciones en tiempo real.",
		Events: []string{
			"order.created",
			"order.updated",
		},
	}, nil
}
