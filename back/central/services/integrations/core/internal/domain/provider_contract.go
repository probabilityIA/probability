package domain

import (
	"context"
	"fmt"
)

// ErrNotSupported indica que la operación no está soportada por esta integración.
var ErrNotSupported = fmt.Errorf("operation not supported by this integration")

// IIntegrationContract es la interfaz que TODA integración debe satisfacer.
// Los providers que no soportan una operación deben embedear BaseIntegration,
// que retorna ErrNotSupported para todos los métodos opcionales.
type IIntegrationContract interface {
	// Obligatorio — toda integración debe implementar esto explícitamente.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error

	// Sincronización de órdenes (ej: Shopify)
	SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error
	SyncOrdersByIntegrationIDWithParams(ctx context.Context, integrationID string, params interface{}) error

	// Webhook — URL informativa
	GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*WebhookInfo, error)

	// Webhook — operaciones CRUD en plataformas externas (ej: Shopify)
	ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error)
	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error
	VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]interface{}, error)
	CreateWebhook(ctx context.Context, integrationID string, baseURL string) (interface{}, error)
}

// BaseIntegration provee implementaciones por defecto que retornan ErrNotSupported.
// Los providers deben embedear este struct y solo sobrescribir los métodos que soportan.
type BaseIntegration struct{}

func (BaseIntegration) TestConnection(_ context.Context, _ map[string]interface{}, _ map[string]interface{}) error {
	return ErrNotSupported
}
func (BaseIntegration) SyncOrdersByIntegrationID(_ context.Context, _ string) error {
	return ErrNotSupported
}
func (BaseIntegration) SyncOrdersByIntegrationIDWithParams(_ context.Context, _ string, _ interface{}) error {
	return ErrNotSupported
}
func (BaseIntegration) GetWebhookURL(_ context.Context, _ string, _ uint) (*WebhookInfo, error) {
	return nil, ErrNotSupported
}
func (BaseIntegration) ListWebhooks(_ context.Context, _ string) ([]interface{}, error) {
	return nil, ErrNotSupported
}
func (BaseIntegration) DeleteWebhook(_ context.Context, _, _ string) error {
	return ErrNotSupported
}
func (BaseIntegration) VerifyWebhooksByURL(_ context.Context, _ string, _ string) ([]interface{}, error) {
	return nil, ErrNotSupported
}
func (BaseIntegration) CreateWebhook(_ context.Context, _ string, _ string) (interface{}, error) {
	return nil, ErrNotSupported
}
