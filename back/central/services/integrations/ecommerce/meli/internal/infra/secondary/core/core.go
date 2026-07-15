package core

import (
	"context"
	"fmt"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases"
)

type MeliCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IMeliUseCase
}

func New(useCase usecases.IMeliUseCase) *MeliCore {
	return &MeliCore{useCase: useCase}
}

func (m *MeliCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return m.useCase.TestConnection(ctx, config, credentials)
}

func (m *MeliCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return m.useCase.SyncOrders(ctx, integrationID)
}

func (m *MeliCore) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	return m.useCase.SyncOrdersWithParams(ctx, integrationID, params)
}

func (m *MeliCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	webhookURL := fmt.Sprintf("%s/meli/notifications", baseURL)

	return &integrationcore.WebhookInfo{
		URL:    webhookURL,
		Method: "POST",
		Description: "Configura esta URL en MercadoLibre Developers > Notificaciones " +
			"para recibir eventos de ordenes, envios, pagos, items y reclamos en tiempo real.",
		Events: []string{
			"orders_v2",
			"shipments",
			"payments",
			"items",
			"claims",
		},
	}, nil
}
