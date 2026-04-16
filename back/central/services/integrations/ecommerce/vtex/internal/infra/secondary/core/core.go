package core

import (
	"context"
	"fmt"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases"
)

// VTEXCore implementa integrationcore.IIntegrationContract para VTEX.
// Embeds BaseIntegration para los métodos opcionales no soportados.
type VTEXCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IVTEXUseCase
}

// New crea el adaptador core de VTEX.
func New(useCase usecases.IVTEXUseCase) *VTEXCore {
	return &VTEXCore{useCase: useCase}
}

// TestConnection delega al use case la verificación de credenciales.
func (v *VTEXCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return v.useCase.TestConnection(ctx, config, credentials)
}

// SyncOrdersByIntegrationID delega la sincronización de órdenes al use case.
func (v *VTEXCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return v.useCase.SyncOrders(ctx, integrationID)
}

// SyncOrdersByIntegrationIDWithParams delega la sincronización con parámetros al use case.
func (v *VTEXCore) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	return v.useCase.SyncOrdersWithParams(ctx, integrationID, params)
}

// GetWebhookURL retorna la URL para los webhooks de VTEX.
func (v *VTEXCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/integrations/vtex/webhook", baseURL)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Configura este webhook en VTEX > Master Data o Orders Feed. " +
			"Suscríbete a los eventos de órdenes para recibir notificaciones en tiempo real.",
		Events: []string{
			"order.created",
			"order.updated",
			"order.cancelled",
		},
	}, nil
}
