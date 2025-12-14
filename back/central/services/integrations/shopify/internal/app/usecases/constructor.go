package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

type SyncOrdersUseCase struct {
	integrationService domain.IIntegrationService
	shopifyClient      domain.ShopifyClient
	orderPublisher     domain.OrderPublisher
}

// IShopifyUseCase define la interfaz para los casos de uso de Shopify
type IShopifyUseCase interface {
	SyncOrders(ctx context.Context, integrationID string) error
	GetOrders(ctx context.Context, integration *domain.Integration, storeDomain, accessToken string, params *domain.GetOrdersParams) error
	GetOrder(ctx context.Context, integrationID string, orderID string) error
	CreateOrder(ctx context.Context, shopDomain string, order *domain.ShopifyOrder) error
	ProcessOrderPaid(ctx context.Context, shopDomain string, order *domain.ShopifyOrder) error
	ProcessOrderUpdated(ctx context.Context, shopDomain string, order *domain.ShopifyOrder) error
	ProcessOrderCancelled(ctx context.Context, shopDomain string, order *domain.ShopifyOrder) error
	ProcessOrderFulfilled(ctx context.Context, shopDomain string, order *domain.ShopifyOrder) error
	ProcessOrderPartiallyFulfilled(ctx context.Context, shopDomain string, order *domain.ShopifyOrder) error
}

// New crea una nueva instancia de IShopifyUseCase
func New(integrationService domain.IIntegrationService, shopifyClient domain.ShopifyClient, orderPublisher domain.OrderPublisher) IShopifyUseCase {
	return &SyncOrdersUseCase{
		integrationService: integrationService,
		shopifyClient:      shopifyClient,
		orderPublisher:     orderPublisher,
	}
}
