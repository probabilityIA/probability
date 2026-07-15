package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

type IMeliClient interface {
	TestConnection(ctx context.Context, accessToken string) error

	GetOrder(ctx context.Context, accessToken string, orderID int64) (*MeliOrder, []byte, error)

	GetOrders(ctx context.Context, accessToken string, sellerID int64, params *GetOrdersParams) (*GetOrdersResult, [][]byte, error)

	GetShipmentDetail(ctx context.Context, accessToken string, shipmentID int64) (*MeliShippingDetail, error)

	GetShipmentOrderIDs(ctx context.Context, accessToken string, shipmentID int64) ([]int64, error)

	GetBillingInfo(ctx context.Context, accessToken string, orderID int64) (*MeliBillingInfo, error)

	GetPack(ctx context.Context, accessToken string, packID int64) (*MeliPack, error)

	GetPaymentOrderID(ctx context.Context, accessToken string, paymentID int64) (int64, error)

	GetClaim(ctx context.Context, accessToken string, claimID int64) (*MeliClaim, error)

	RefreshToken(ctx context.Context, appID, clientSecret, refreshToken string) (*TokenResponse, error)

	GetUserMe(ctx context.Context, accessToken string) (*MeliSeller, error)

	GetProducts(ctx context.Context, accessToken string, sellerID int64) ([]MeliProduct, error)

	CreateProduct(ctx context.Context, accessToken string, input CreateProductInput) (string, error)

	GetItem(ctx context.Context, accessToken, itemID string) (*MeliItemDetail, error)

	UpdateStock(ctx context.Context, accessToken, itemID string, quantity int) error

	GetUserProductStock(ctx context.Context, accessToken, userProductID string) (*UserProductStock, error)

	UpdateUserProductStock(ctx context.Context, accessToken, userProductID, version string, locations []StockLocation) error

	SendShipmentStatus(ctx context.Context, accessToken string, shipmentID int64, status string) error
}

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
	GetIntegrationByStoreID(ctx context.Context, storeID string) (*Integration, error)
}

type OrderPublisher interface {
	Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
}
