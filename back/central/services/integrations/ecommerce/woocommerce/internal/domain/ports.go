package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

type IWooCommerceClient interface {
	TestConnection(ctx context.Context, storeURL, consumerKey, consumerSecret string) error

	GetOrders(ctx context.Context, storeURL, consumerKey, consumerSecret string, params *GetOrdersParams) (*GetOrdersResult, [][]byte, error)

	GetOrder(ctx context.Context, storeURL, consumerKey, consumerSecret string, orderID int64) (*WooCommerceOrder, []byte, error)

	CreateWebhook(ctx context.Context, storeURL, consumerKey, consumerSecret, deliveryURL, secret, topic string) (int64, error)

	ListWebhooks(ctx context.Context, storeURL, consumerKey, consumerSecret string) ([]WebhookItem, error)

	DeleteWebhook(ctx context.Context, storeURL, consumerKey, consumerSecret, webhookID string) error

	UpdateProductStock(ctx context.Context, storeURL, consumerKey, consumerSecret, productExternalID string, quantity int) error

	CreateProduct(ctx context.Context, storeURL, consumerKey, consumerSecret string, input CreateProductInput) (string, error)

	GetProducts(ctx context.Context, storeURL, consumerKey, consumerSecret string) ([]WooProduct, error)

	GetProductsStock(ctx context.Context, storeURL, consumerKey, consumerSecret string, externalIDs []string) ([]ChannelStock, error)
}

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
}

type OrderPublisher interface {
	Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
}
