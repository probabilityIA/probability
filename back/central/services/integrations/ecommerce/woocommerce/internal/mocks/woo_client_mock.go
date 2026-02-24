package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// WooClientMock mock de domain.IWooCommerceClient para tests unitarios.
// Cada método de la interfaz tiene su correspondiente campo Fn que permite
// inyectar el comportamiento deseado en cada test.
type WooClientMock struct {
	TestConnectionFn func(ctx context.Context, storeURL, consumerKey, consumerSecret string) error
	GetOrdersFn      func(ctx context.Context, storeURL, consumerKey, consumerSecret string, params *domain.GetOrdersParams) (*domain.GetOrdersResult, [][]byte, error)
	GetOrderFn       func(ctx context.Context, storeURL, consumerKey, consumerSecret string, orderID int64) (*domain.WooCommerceOrder, []byte, error)
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
