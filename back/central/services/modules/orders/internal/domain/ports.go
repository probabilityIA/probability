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
	CreateOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, id string) (*Order, error)
	GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*Order, error)
	ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]Order, int64, error)
	UpdateOrder(ctx context.Context, order *Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderRaw(ctx context.Context, id string) (*OrderChannelMetadata, error)

	// Validation
	OrderExists(ctx context.Context, externalID string, integrationID uint) (bool, error)

	// ============================================
	// MÉTODOS PARA TABLAS RELACIONADAS
	// ============================================

	// OrderItems
	CreateOrderItems(ctx context.Context, items []*OrderItem) error

	// Addresses
	CreateAddresses(ctx context.Context, addresses []*Address) error

	// Payments
	CreatePayments(ctx context.Context, payments []*Payment) error

	// Shipments
	CreateShipments(ctx context.Context, shipments []*Shipment) error

	// ChannelMetadata
	CreateChannelMetadata(ctx context.Context, metadata *OrderChannelMetadata) error

	// ============================================
	// MÉTODOS DE CATÁLOGO (VALIDACIÓN)
	// ============================================

	// Products
	GetProductBySKU(ctx context.Context, businessID uint, sku string) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error

	// Clients
	GetClientByEmail(ctx context.Context, businessID uint, email string) (*Client, error)
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
