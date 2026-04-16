package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

// IVTEXClient define las operaciones del cliente HTTP de VTEX.
// Implementado en infra/secondary/client.
type IVTEXClient interface {
	// TestConnection verifica que las credenciales sean válidas.
	TestConnection(ctx context.Context, storeURL, apiKey, apiToken string) error

	// GetOrders obtiene la lista de órdenes con paginación.
	// Retorna el resumen de órdenes (sin detalle completo).
	GetOrders(ctx context.Context, storeURL, apiKey, apiToken string, page, perPage int, filters map[string]string) (*VTEXOrderListResponse, error)

	// GetOrderByID obtiene el detalle completo de una orden.
	// Retorna la orden tipada y los bytes crudos (para ChannelMetadata.RawData).
	GetOrderByID(ctx context.Context, storeURL, apiKey, apiToken string, orderID string) (*VTEXOrder, []byte, error)
}

// IIntegrationService define las operaciones del core de integraciones
// que el módulo de VTEX necesita.
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
