package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

type IJumpsellerClient interface {
	TestConnection(ctx context.Context, cred Credential) error

	GetStoreInfo(ctx context.Context, cred Credential) (*StoreInfo, error)

	GetOrders(ctx context.Context, cred Credential, params *GetOrdersParams) (*GetOrdersResult, [][]byte, error)

	GetOrder(ctx context.Context, cred Credential, orderID int64) (*JumpsellerOrder, []byte, error)

	UpdateOrder(ctx context.Context, cred Credential, orderID int64, fields UpdateOrderFields) error

	CreateHook(ctx context.Context, cred Credential, event, url string) (string, error)

	ListHooks(ctx context.Context, cred Credential) ([]WebhookItem, error)

	DeleteHook(ctx context.Context, cred Credential, hookID string) error

	GetProducts(ctx context.Context, cred Credential) ([]JumpsellerProduct, error)

	GetLocations(ctx context.Context, cred Credential) ([]Location, error)

	ResolveStockTarget(ctx context.Context, cred Credential, sku string) (*StockTarget, error)

	SetProductStock(ctx context.Context, cred Credential, productID int64, stock int) error

	SetVariantStock(ctx context.Context, cred Credential, productID, variantID int64, stock int) error

	CreateProduct(ctx context.Context, cred Credential, input CreateProductInput) (string, error)

	UpdateProduct(ctx context.Context, cred Credential, productID int64, input UpdateProductInput) error

	RefreshToken(ctx context.Context, tokenURL, clientID, clientSecret, refreshToken string) (*TokenResponse, error)
}

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
	UpdateIntegrationCredentials(ctx context.Context, integrationID string, credentials map[string]interface{}) error
	GetPlatformCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
}

type OrderPublisher interface {
	Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
}
