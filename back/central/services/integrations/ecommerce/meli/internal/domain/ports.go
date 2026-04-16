package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

// IMeliClient define las operaciones del cliente HTTP de MercadoLibre.
// Implementado en infra/secondary/client.
type IMeliClient interface {
	// TestConnection verifica que las credenciales sean válidas.
	TestConnection(ctx context.Context, accessToken string) error

	// GetOrder obtiene una orden específica por ID.
	// Retorna la orden tipada, los bytes crudos (para RawData), y error.
	GetOrder(ctx context.Context, accessToken string, orderID int64) (*MeliOrder, []byte, error)

	// GetOrders obtiene órdenes del vendedor con paginación.
	// Retorna el resultado paginado, los bytes crudos por orden, y error.
	GetOrders(ctx context.Context, accessToken string, sellerID int64, params *GetOrdersParams) (*GetOrdersResult, [][]byte, error)

	// GetShipmentDetail obtiene los detalles completos de un envío.
	// La orden solo trae shipping.id, se necesita esta llamada para dirección y estado.
	GetShipmentDetail(ctx context.Context, accessToken string, shipmentID int64) (*MeliShippingDetail, error)

	// RefreshToken renueva el access_token usando el refresh_token.
	RefreshToken(ctx context.Context, appID, clientSecret, refreshToken string) (*TokenResponse, error)

	// GetUserMe obtiene los datos del usuario autenticado (para extraer seller_id).
	GetUserMe(ctx context.Context, accessToken string) (*MeliSeller, error)
}

// IIntegrationService define las operaciones del core de integraciones
// que el módulo de MercadoLibre necesita.
type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
	// GetIntegrationByStoreID busca integración por store_id (seller_id de MeLi).
	GetIntegrationByStoreID(ctx context.Context, storeID string) (*Integration, error)
}

// OrderPublisher publica órdenes al canal canónico de RabbitMQ.
// Implementado en infra/secondary/queue.
type OrderPublisher interface {
	Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
}
