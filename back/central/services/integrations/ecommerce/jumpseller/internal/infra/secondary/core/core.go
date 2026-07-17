package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases"
)

type JumpsellerCore struct {
	integrationcore.BaseIntegration
	useCase usecases.IJumpsellerUseCase
}

func New(useCase usecases.IJumpsellerUseCase) *JumpsellerCore {
	return &JumpsellerCore{useCase: useCase}
}

func (j *JumpsellerCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return j.useCase.TestConnection(ctx, config, credentials)
}

func (j *JumpsellerCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return j.useCase.SyncOrders(ctx, integrationID)
}

func (j *JumpsellerCore) SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error {
	return j.useCase.SyncOrdersWithParams(ctx, integrationID, params)
}

func (j *JumpsellerCore) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	return j.useCase.UpdateInventory(ctx, integrationID, productExternalID, quantity)
}

func (j *JumpsellerCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*integrationcore.WebhookInfo, error) {
	return &integrationcore.WebhookInfo{
		URL:    usecases.WebhookDeliveryURL(baseURL, integrationID),
		Method: "POST",
		Description: "Probability crea estos webhooks automaticamente en tu tienda Jumpseller. " +
			"Recibiras las ordenes en tiempo real cuando se creen, se paguen, se envien o se cancelen.",
		Events: []string{
			"order_created",
			"order_paid",
			"order_shipped",
			"order_canceled",
		},
	}, nil
}

func (j *JumpsellerCore) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	items, err := j.useCase.ListWebhooks(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	return result, nil
}

func (j *JumpsellerCore) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	return j.useCase.DeleteWebhook(ctx, integrationID, webhookID)
}

func (j *JumpsellerCore) VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]interface{}, error) {
	return j.ListWebhooks(ctx, integrationID)
}

func (j *JumpsellerCore) CreateWebhook(ctx context.Context, integrationID string, baseURL string) (interface{}, error) {
	if err := j.useCase.CreateWebhooks(ctx, integrationID, baseURL); err != nil {
		return nil, err
	}

	items, err := j.useCase.ListWebhooks(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	created := make([]string, 0, len(items))
	webhookURL := ""
	for _, item := range items {
		created = append(created, item.Topic)
		webhookURL = item.Address
	}

	return map[string]interface{}{
		"WebhookURL":       webhookURL,
		"CreatedWebhooks":  created,
		"ExistingWebhooks": []interface{}{},
		"DeletedWebhooks":  []interface{}{},
	}, nil
}
