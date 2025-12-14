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
}

type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	GetIntegrationByStoreID(ctx context.Context, storeID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
}
