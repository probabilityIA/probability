package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

type ITiendanubeClient interface {
	TestConnection(ctx context.Context, cred Credential) error
	GetStoreInfo(ctx context.Context, cred Credential) (*StoreInfo, error)
	GetProducts(ctx context.Context, cred Credential) ([]TiendanubeProduct, error)
	ResolveStockTarget(ctx context.Context, cred Credential, sku string) (*StockTarget, error)
	SetVariantStock(ctx context.Context, cred Credential, productID, variantID int64, stock int) error
	GetProductsStock(ctx context.Context, cred Credential, externalIDs []string) ([]ChannelStock, error)
	CreateProduct(ctx context.Context, cred Credential, input CreateProductInput) (int64, int64, error)
	UpdateProduct(ctx context.Context, cred Credential, productID int64, input UpdateProductInput) error
	UpdateVariant(ctx context.Context, cred Credential, productID, variantID int64, input UpdateVariantInput) error
	ListWebhooks(ctx context.Context, cred Credential) ([]WebhookItem, error)
	CreateWebhook(ctx context.Context, cred Credential, event, webhookURL string) (string, error)
	DeleteWebhook(ctx context.Context, cred Credential, webhookID string) error
	GetOrder(ctx context.Context, cred Credential, orderID string) (*TiendanubeOrder, []byte, error)
	GetOrders(ctx context.Context, cred Credential, filters OrderFilters) ([]TiendanubeOrder, error)
	ListFulfillmentOrders(ctx context.Context, cred Credential, orderID string) ([]FulfillmentOrder, error)
	UpdateFulfillmentOrder(ctx context.Context, cred Credential, orderID, fulfillmentOrderID, status string, tracking *TrackingInfo) error
	CreateTrackingEvent(ctx context.Context, cred Credential, orderID, fulfillmentOrderID string, event TrackingEvent) error
	CancelOrder(ctx context.Context, cred Credential, orderID string) error
}

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
}

type OrderPublisher interface {
	Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
}
