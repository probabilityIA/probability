package domain

import (
	"context"
)

// ───────────────────────────────────────────
//
//	REPOSITORY INTERFACE
//
// ───────────────────────────────────────────

// IRepository define todos los métodos de repositorio del módulo orders
type IRepository interface {
	// CRUD Operations
	CreateOrder(ctx context.Context, order *ProbabilityOrder) error
	GetOrderByID(ctx context.Context, id string) (*ProbabilityOrder, error)
	GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*ProbabilityOrder, error)
	ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]ProbabilityOrder, int64, error)
	UpdateOrder(ctx context.Context, order *ProbabilityOrder) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderRaw(ctx context.Context, id string) (*ProbabilityOrderChannelMetadata, error)
	CountOrdersByClientID(ctx context.Context, clientID uint) (int64, error)

	// Validation
	OrderExists(ctx context.Context, externalID string, integrationID uint) (bool, error)
	GetOrderByExternalID(ctx context.Context, externalID string, integrationID uint) (*ProbabilityOrder, error)

	// ============================================
	// MÉTODOS PARA TABLAS RELACIONADAS
	// ============================================

	// OrderItems
	CreateOrderItems(ctx context.Context, items []*ProbabilityOrderItem) error

	// Addresses
	CreateAddresses(ctx context.Context, addresses []*ProbabilityAddress) error

	// Payments
	CreatePayments(ctx context.Context, payments []*ProbabilityPayment) error

	// Shipments
	CreateShipments(ctx context.Context, shipments []*ProbabilityShipment) error

	// ChannelMetadata
	CreateChannelMetadata(ctx context.Context, metadata *ProbabilityOrderChannelMetadata) error

	// ============================================
	// MÉTODOS DE CATÁLOGO (VALIDACIÓN)
	// ============================================

	// Products
	GetProductBySKU(ctx context.Context, businessID uint, sku string) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error

	// Clients
	GetClientByEmail(ctx context.Context, businessID uint, email string) (*Client, error)
	GetClientByDNI(ctx context.Context, businessID uint, dni string) (*Client, error)
	CreateClient(ctx context.Context, client *Client) error
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
	MapAndSaveOrder(ctx context.Context, dto *ProbabilityOrderDTO) (*OrderResponse, error)
	UpdateOrder(ctx context.Context, existingOrder *ProbabilityOrder, dto *ProbabilityOrderDTO) (*OrderResponse, error)
}

// ───────────────────────────────────────────
//
//	EVENT PUBLISHER INTERFACE
//
// ───────────────────────────────────────────

// IOrderEventPublisher define la interfaz para publicar eventos de órdenes
type IOrderEventPublisher interface {
	// PublishOrderEvent publica un evento de orden a Redis
	PublishOrderEvent(ctx context.Context, event *OrderEvent) error
}

// IOrderScoreUseCase define la interfaz para el caso de uso de cálculo de score
type IOrderScoreUseCase interface {
	CalculateOrderScore(order *ProbabilityOrder) (float64, []string)
	CalculateAndUpdateOrderScore(ctx context.Context, orderID string) error
}
