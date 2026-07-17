package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

type VTEXCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IVTEXUseCase
}

func New(useCase usecases.IVTEXUseCase) *VTEXCore {
	return &VTEXCore{useCase: useCase}
}

func (v *VTEXCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return v.useCase.TestConnection(ctx, config, credentials)
}

func (v *VTEXCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return v.useCase.SyncOrders(ctx, integrationID)
}

func (v *VTEXCore) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	return v.useCase.SyncOrdersWithParams(ctx, integrationID, params)
}

func (v *VTEXCore) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	return v.useCase.UpdateInventory(ctx, integrationID, productExternalID, quantity)
}

func (v *VTEXCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	return &integrationcore.WebhookInfo{
		URL:    usecases.WebhookDeliveryURL(baseURL, integrationID),
		Method: "POST",
		Description: "VTEX admite una sola configuracion de hook por cuenta: al registrarla se reemplaza " +
			"cualquier hook anterior, incluido el de otra herramienta. Probability consulta el hook actual " +
			"antes de registrar y avisa si encuentra uno ajeno.",
		Events: domain.WebhookOrderStates,
	}, nil
}

func (v *VTEXCore) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	items, err := v.useCase.ListWebhooks(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":       item.ID,
			"address":  item.Address,
			"statuses": item.Statuses,
			"is_ours":  item.IsOurs,
		})
	}
	return result, nil
}

func (v *VTEXCore) VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]interface{}, error) {
	item, err := v.useCase.InspectWebhook(ctx, integrationID, baseURL)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return []interface{}{}, nil
	}
	return []interface{}{map[string]interface{}{
		"id":       item.ID,
		"address":  item.Address,
		"statuses": item.Statuses,
		"is_ours":  item.IsOurs,
	}}, nil
}

func (v *VTEXCore) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	return v.useCase.DeleteWebhook(ctx, integrationID, webhookID)
}

func (v *VTEXCore) CreateWebhook(ctx context.Context, integrationID string, baseURL string) (interface{}, error) {
	if err := v.useCase.CreateWebhooks(ctx, integrationID, baseURL, false); err != nil {
		return nil, err
	}

	item, err := v.useCase.InspectWebhook(ctx, integrationID, baseURL)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"webhook_url": usecases.WebhookDeliveryURL(baseURL, 0),
		"events":      domain.WebhookOrderStates,
	}
	if item != nil {
		result["webhook_url"] = item.Address
		result["is_ours"] = item.IsOurs
	}
	return result, nil
}
