package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// ───────────────────────────────────────────
//
//	REPOSITORY INTERFACE
//
// ───────────────────────────────────────────

// IRepository define todos los métodos de repositorio del módulo orders
type IRepository interface {
	// CRUD Operations
	CreateOrder(ctx context.Context, order *entities.ProbabilityOrder) error
	GetOrderByID(ctx context.Context, id string) (*entities.ProbabilityOrder, error)
	GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*entities.ProbabilityOrder, error)
	GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*entities.ProbabilityOrder, error)
	ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.ProbabilityOrder, int64, error)
	UpdateOrder(ctx context.Context, order *entities.ProbabilityOrder) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderRaw(ctx context.Context, id string) (*entities.ProbabilityOrderChannelMetadata, error)
	CountOrdersByClientID(ctx context.Context, clientID uint) (int64, error)
	GetLastManualOrderNumber(ctx context.Context, businessID uint) (int, error)
	GetFirstIntegrationIDByBusinessID(ctx context.Context, businessID uint) (uint, error)
	GetPlatformIntegrationIDByBusinessID(ctx context.Context, businessID uint) (uint, error)

	// Validation
	OrderExists(ctx context.Context, externalID string, integrationID uint) (bool, error)
	GetOrderByExternalID(ctx context.Context, externalID string, integrationID uint) (*entities.ProbabilityOrder, error)

	// ============================================
	// MÉTODOS PARA TABLAS RELACIONADAS
	// ============================================

	// OrderItems
	CreateOrderItems(ctx context.Context, items []*entities.ProbabilityOrderItem) error

	// Addresses
	CreateAddresses(ctx context.Context, addresses []*entities.ProbabilityAddress) error

	// Payments
	CreatePayments(ctx context.Context, payments []*entities.ProbabilityPayment) error

	// Shipments
	CreateShipments(ctx context.Context, shipments []*entities.ProbabilityShipment) error

	// ChannelMetadata
	CreateChannelMetadata(ctx context.Context, metadata *entities.ProbabilityOrderChannelMetadata) error

	// ============================================
	// MÉTODOS DE CATÁLOGO (VALIDACIÓN)
	// ============================================

	// Products
	GetProductBySKU(ctx context.Context, businessID uint, sku string) (*entities.Product, error)
	CreateProduct(ctx context.Context, product *entities.Product) error

	// Clients
	GetClientByEmail(ctx context.Context, businessID uint, email string) (*entities.Client, error)
	GetClientByDNI(ctx context.Context, businessID uint, dni string) (*entities.Client, error)
	CreateClient(ctx context.Context, client *entities.Client) error

	// OrderErrors
	CreateOrderError(ctx context.Context, orderError *entities.OrderError) error

	// ============================================
	// MÉTODOS DE CONSULTA A TABLAS DE ESTADOS
	// (Replicados localmente - no compartir repos)
	// ============================================

	// OrderStatuses - Consultas a la tabla order_statuses
	GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error)

	// PaymentStatuses - Consultas a la tabla payment_statuses
	GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error)

	// FulfillmentStatuses - Consultas a la tabla fulfillment_statuses
	GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error)
}

// ───────────────────────────────────────────
//
//	ORDER CONSUMER INTERFACE
//
// ───────────────────────────────────────────

// IOrderConsumer define la interfaz para consumir órdenes desde colas
type IOrderConsumer interface {
	// Start inicia el consumidor de órdenes
	Start(ctx context.Context) error
}

// ───────────────────────────────────────────
//
//	USE CASE INTERFACES
//
// ───────────────────────────────────────────

// IOrderMappingUseCase define el caso de uso para mapear y guardar órdenes desde integraciones
type IOrderMappingUseCase interface {
	MapAndSaveOrder(ctx context.Context, dto *dtos.ProbabilityOrderDTO) (*dtos.OrderResponse, error)
	UpdateOrder(ctx context.Context, existingOrder *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) (*dtos.OrderResponse, error)
}

// ───────────────────────────────────────────
//
//	EVENT PUBLISHER INTERFACE
//
// ───────────────────────────────────────────

// IOrderEventPublisher define la interfaz para publicar eventos de órdenes
type IOrderEventPublisher interface {
	// PublishOrderEvent publica un evento de orden a Redis con snapshot completo
	// El parámetro order permite construir el OrderSnapshot completo sin queries adicionales
	PublishOrderEvent(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error
}

// IOrderScoreUseCase define la interfaz para el caso de uso de cálculo de score
type IOrderScoreUseCase interface {
	CalculateOrderScore(order *entities.ProbabilityOrder) (float64, []string)
	CalculateAndUpdateOrderScore(ctx context.Context, orderID string) error
}

// IRequestConfirmationUseCase define la interfaz del caso de uso de solicitud de confirmación
type IRequestConfirmationUseCase interface {
	RequestConfirmation(ctx context.Context, orderID string) error
}

// IOrderUseCase define la interfaz para los casos de uso CRUD de órdenes
type IOrderUseCase interface {
	CreateOrder(ctx context.Context, req *dtos.CreateOrderRequest) (*dtos.OrderResponse, error)
	GetOrderByID(ctx context.Context, id string) (*dtos.OrderResponse, error)
	GetOrderRaw(ctx context.Context, id string) (*dtos.OrderRawResponse, error)
	ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*dtos.OrdersListResponse, error)
	UpdateOrder(ctx context.Context, id string, req *dtos.UpdateOrderRequest) (*dtos.OrderResponse, error)
	DeleteOrder(ctx context.Context, id string) error
}

// ───────────────────────────────────────────
//
//	RABBITMQ PUBLISHER INTERFACE
//
// ───────────────────────────────────────────

// IIntegrationEventPublisher publica eventos de sincronización de integraciones al exchange de eventos
type IIntegrationEventPublisher interface {
	PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, orderNumber, externalID, platform, reason, errMsg string)
	PublishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
	PublishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
}

// IOrderRabbitPublisher publica eventos de órdenes a RabbitMQ
type IOrderRabbitPublisher interface {
	// Eventos de ciclo de vida
	PublishOrderCreated(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderUpdated(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderCancelled(ctx context.Context, order *entities.ProbabilityOrder) error

	// Eventos de estado
	PublishOrderStatusChanged(ctx context.Context, order *entities.ProbabilityOrder, previousStatus, currentStatus string) error

	// Eventos de confirmación (ya existe)
	PublishConfirmationRequested(ctx context.Context, order *entities.ProbabilityOrder) error

	// Método genérico para cualquier evento con snapshot completo
	PublishOrderEvent(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error
}
