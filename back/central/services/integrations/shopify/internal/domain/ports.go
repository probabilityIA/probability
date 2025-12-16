package domain

import (
	"context"
)

type OrderPublisher interface {
	Publish(ctx context.Context, order *ProbabilityOrderDTO) error
}

type ShopifyClient interface {
	ValidateToken(ctx context.Context, storeName, accessToken string) (bool, map[string]interface{}, error)
	GetOrders(ctx context.Context, storeName, accessToken string, params *GetOrdersParams) ([]ShopifyOrder, string, error)
	GetOrder(ctx context.Context, storeName, accessToken string, orderID string) (*ShopifyOrder, error)
	CreateWebhook(ctx context.Context, storeName, accessToken, webhookURL, event string) (string, error) // Crea un webhook en Shopify y retorna el webhook ID
	ListWebhooks(ctx context.Context, storeName, accessToken string) ([]WebhookInfo, error)              // Lista todos los webhooks de la tienda
	DeleteWebhook(ctx context.Context, storeName, accessToken, webhookID string) error                   // Elimina un webhook por ID
	SetDebug(enabled bool)                                                                               // Habilita logging de peticiones HTTP
}

// WebhookInfo representa la información de un webhook de Shopify
type WebhookInfo struct {
	ID        string `json:"id"`
	Address   string `json:"address"`
	Topic     string `json:"topic"`
	Format    string `json:"format"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	GetIntegrationByStoreID(ctx context.Context, storeID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error // Actualiza el config de una integración
}
