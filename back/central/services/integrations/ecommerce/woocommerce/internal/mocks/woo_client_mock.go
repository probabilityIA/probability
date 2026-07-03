package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// WooClientMock mock de domain.IWooCommerceClient para tests unitarios.
// Cada método de la interfaz tiene su correspondiente campo Fn que permite
// inyectar el comportamiento deseado en cada test.
type WooClientMock struct {
	TestConnectionFn     func(ctx context.Context, storeURL, consumerKey, consumerSecret string) error
	GetOrdersFn          func(ctx context.Context, storeURL, consumerKey, consumerSecret string, params *domain.GetOrdersParams) (*domain.GetOrdersResult, [][]byte, error)
	GetOrderFn           func(ctx context.Context, storeURL, consumerKey, consumerSecret string, orderID int64) (*domain.WooCommerceOrder, []byte, error)
	CreateWebhookFn      func(ctx context.Context, storeURL, consumerKey, consumerSecret, deliveryURL, secret, topic string) (int64, error)
	UpdateProductStockFn func(ctx context.Context, storeURL, consumerKey, consumerSecret, productExternalID string, quantity int) error
	CreateProductFn      func(ctx context.Context, storeURL, consumerKey, consumerSecret string, input domain.CreateProductInput) (string, error)
}

func (m *WooClientMock) CreateProduct(ctx context.Context, storeURL, consumerKey, consumerSecret string, input domain.CreateProductInput) (string, error) {
	if m.CreateProductFn != nil {
		return m.CreateProductFn(ctx, storeURL, consumerKey, consumerSecret, input)
	}
	return "0", nil
}

// Verificar en tiempo de compilación que WooClientMock implementa la interfaz.
var _ domain.IWooCommerceClient = (*WooClientMock)(nil)

func (m *WooClientMock) TestConnection(ctx context.Context, storeURL, consumerKey, consumerSecret string) error {
	if m.TestConnectionFn != nil {
		return m.TestConnectionFn(ctx, storeURL, consumerKey, consumerSecret)
	}
	return nil
}

func (m *WooClientMock) GetOrders(ctx context.Context, storeURL, consumerKey, consumerSecret string, params *domain.GetOrdersParams) (*domain.GetOrdersResult, [][]byte, error) {
	if m.GetOrdersFn != nil {
		return m.GetOrdersFn(ctx, storeURL, consumerKey, consumerSecret, params)
	}
	return &domain.GetOrdersResult{}, nil, nil
}

func (m *WooClientMock) GetOrder(ctx context.Context, storeURL, consumerKey, consumerSecret string, orderID int64) (*domain.WooCommerceOrder, []byte, error) {
	if m.GetOrderFn != nil {
		return m.GetOrderFn(ctx, storeURL, consumerKey, consumerSecret, orderID)
	}
	return &domain.WooCommerceOrder{}, nil, nil
}

func (m *WooClientMock) CreateWebhook(ctx context.Context, storeURL, consumerKey, consumerSecret, deliveryURL, secret, topic string) (int64, error) {
	if m.CreateWebhookFn != nil {
		return m.CreateWebhookFn(ctx, storeURL, consumerKey, consumerSecret, deliveryURL, secret, topic)
	}
	return 0, nil
}

func (m *WooClientMock) UpdateProductStock(ctx context.Context, storeURL, consumerKey, consumerSecret, productExternalID string, quantity int) error {
	if m.UpdateProductStockFn != nil {
		return m.UpdateProductStockFn(ctx, storeURL, consumerKey, consumerSecret, productExternalID, quantity)
	}
	return nil
}
