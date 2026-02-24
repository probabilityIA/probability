package core

import (
	"context"
	"fmt"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/amazon/internal/app/usecases"
)

// AmazonCore implementa integrationcore.IIntegrationContract para Amazon.
// Embeds BaseIntegration para los metodos opcionales no soportados.
type AmazonCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IAmazonUseCase
}

// New crea el adaptador core de Amazon.
func New(useCase usecases.IAmazonUseCase) *AmazonCore {
	return &AmazonCore{useCase: useCase}
}

// TestConnection delega al use case la verificacion de credenciales.
func (a *AmazonCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return a.useCase.TestConnection(ctx, config, credentials)
}

// GetWebhookURL retorna la URL para las notificaciones de Amazon SP-API.
func (a *AmazonCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/integrations/amazon/notification", baseURL)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Amazon SP-API usa SQS/SNS para notificaciones. " +
			"Configura la suscripcion en Seller Central > Notifications para recibir eventos.",
		Events: []string{
			"ORDER_CHANGE",
			"LISTINGS_ITEM_STATUS_CHANGE",
		},
	}, nil
}
