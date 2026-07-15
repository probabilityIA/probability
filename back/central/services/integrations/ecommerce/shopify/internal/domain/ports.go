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
	GetOrdersByURL(ctx context.Context, nextPageURL, accessToken string) ([]ShopifyOrder, string, error)
	GetOrder(ctx context.Context, storeName, accessToken string, orderID string) (*ShopifyOrder, error)
	CreateWebhook(ctx context.Context, storeName, accessToken, webhookURL, event string) (string, error)
	ListWebhooks(ctx context.Context, storeName, accessToken string) ([]WebhookInfo, error)
	DeleteWebhook(ctx context.Context, storeName, accessToken, webhookID string) error
	CreateCarrierService(ctx context.Context, storeName, accessToken, callbackURL, name string) (string, error)
	DeleteCarrierService(ctx context.Context, storeName, accessToken, carrierServiceID string) error
	GetLocations(ctx context.Context, storeName, accessToken string) ([]ShopifyLocation, error)
	GetProduct(ctx context.Context, storeName, accessToken, productID string) (*ShopifyProduct, error)
	SetInventoryLevel(ctx context.Context, storeName, accessToken string, locationID, inventoryItemID int64, available int) error
	SetDebug(enabled bool)
}

type ISyncEventPublisher interface {
	PublishSyncEvent(ctx context.Context, integrationID uint, businessID *uint, eventType string, data map[string]interface{})
}

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
}
