package core

import (
	"context"
	"fmt"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases"
)

// MeliCore implementa integrationcore.IIntegrationContract para MercadoLibre.
// Embeds BaseIntegration para los métodos opcionales no soportados.
type MeliCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IMeliUseCase
}

// New crea el adaptador core de MercadoLibre.
func New(useCase usecases.IMeliUseCase) *MeliCore {
	return &MeliCore{useCase: useCase}
}

// TestConnection delega al use case la verificación de credenciales.
func (m *MeliCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return m.useCase.TestConnection(ctx, config, credentials)
}

// SyncOrdersByIntegrationID delega la sincronización de órdenes al use case.
func (m *MeliCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return m.useCase.SyncOrders(ctx, integrationID)
}

// SyncOrdersByIntegrationIDWithParams delega la sincronización con parámetros al use case.
func (m *MeliCore) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	return m.useCase.SyncOrdersWithParams(ctx, integrationID, params)
}

// GetWebhookURL retorna la URL base para los webhooks de MercadoLibre.
// MercadoLibre usa notificaciones (IPN) en lugar de webhooks clásicos.
func (m *MeliCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/integrations/meli/notifications", baseURL)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Configura esta URL en MercadoLibre Developers > Mis aplicaciones > " +
			"Notificaciones para recibir eventos de órdenes y pagos en tiempo real.",
		Events: []string{
			"orders_v2",
			"payments",
			"items",
		},
	}, nil
}
