package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

// IWooCommerceClient define las operaciones del cliente HTTP de WooCommerce.
// Implementado en infra/secondary/client.
type IWooCommerceClient interface {
	// TestConnection verifica que las credenciales sean válidas
	TestConnection(ctx context.Context, storeURL, consumerKey, consumerSecret string) error

	// GetOrders obtiene órdenes paginadas de la tienda WooCommerce.
	// Retorna las órdenes tipadas, los bytes crudos por orden (para RawData), y error.
	GetOrders(ctx context.Context, storeURL, consumerKey, consumerSecret string, params *GetOrdersParams) (*GetOrdersResult, [][]byte, error)

	// GetOrder obtiene una orden específica por ID.
	// Retorna la orden tipada, los bytes crudos, y error.
	GetOrder(ctx context.Context, storeURL, consumerKey, consumerSecret string, orderID int64) (*WooCommerceOrder, []byte, error)
}

// IIntegrationService define las operaciones del core de integraciones
// que el módulo de WooCommerce necesita.
type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
}

// OrderPublisher publica órdenes al canal canónico de RabbitMQ.
// Implementado en infra/secondary/queue.
type OrderPublisher interface {
	Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
}
