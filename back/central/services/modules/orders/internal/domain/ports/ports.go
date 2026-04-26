package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

type IRepository interface {
	CreateOrder(ctx context.Context, order *entities.ProbabilityOrder) error
	GetOrderByID(ctx context.Context, id string) (*entities.ProbabilityOrder, error)
	GetOrderByInternalNumber(ctx context.Context, internalNumber string) (*entities.ProbabilityOrder, error)
	GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*entities.ProbabilityOrder, error)
	GetOrderByOrderNumberAndBusiness(ctx context.Context, orderNumber string, businessID uint) (*entities.ProbabilityOrder, error)
	ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.ProbabilityOrder, int64, error)
	UpdateOrder(ctx context.Context, order *entities.ProbabilityOrder) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderRaw(ctx context.Context, id string) (*entities.ProbabilityOrderChannelMetadata, error)
	CountOrdersByClientID(ctx context.Context, clientID uint) (int64, error)
	GetLastManualOrderNumber(ctx context.Context, businessID uint) (int, error)
	GetBusinessOrderPrefix(ctx context.Context, businessID uint) (string, error)
	GetFirstIntegrationIDByBusinessID(ctx context.Context, businessID uint) (uint, error)
	GetPlatformIntegrationIDByBusinessID(ctx context.Context, businessID uint) (uint, error)

	OrderExists(ctx context.Context, externalID string, integrationID uint) (bool, error)
	GetOrderByExternalID(ctx context.Context, externalID string, integrationID uint) (*entities.ProbabilityOrder, error)

	CreateOrderItems(ctx context.Context, items []*entities.ProbabilityOrderItem) error
	CreateAddresses(ctx context.Context, addresses []*entities.ProbabilityAddress) error
	CreatePayments(ctx context.Context, payments []*entities.ProbabilityPayment) error
	CreateShipments(ctx context.Context, shipments []*entities.ProbabilityShipment) error
	CreateChannelMetadata(ctx context.Context, metadata *entities.ProbabilityOrderChannelMetadata) error

	GetProductBySKU(ctx context.Context, businessID uint, sku string) (*entities.Product, error)
	ResolveProductForOrderItem(ctx context.Context, businessID uint, integrationID uint, item dtos.ProbabilityOrderItemDTO) (*entities.Product, error)
	CreateProduct(ctx context.Context, product *entities.Product) error
	UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID uint, integrationID uint, item dtos.ProbabilityOrderItemDTO) error
	UpdateProductPrice(ctx context.Context, productID string, price float64) error

	GetClientByEmail(ctx context.Context, businessID uint, email string) (*entities.Client, error)
	GetClientByDNI(ctx context.Context, businessID uint, dni string) (*entities.Client, error)
	CreateClient(ctx context.Context, client *entities.Client) error

	CreateOrderError(ctx context.Context, orderError *entities.OrderError) error

	GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error)
	GetOrderStatusIDByCode(ctx context.Context, code string) (*uint, error)
	GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error)
	GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error)

	UpdateOrderStatus(ctx context.Context, orderID string, status string, statusID *uint) error

	CreateOrderHistory(ctx context.Context, history *entities.OrderHistory) error
	GetOrderHistory(ctx context.Context, orderID string) ([]entities.OrderHistory, error)
}

type IOrderConsumer interface {
	Start(ctx context.Context) error
}

type IOrderCreateUseCase interface {
	MapAndSaveOrder(ctx context.Context, dto *dtos.ProbabilityOrderDTO) (*dtos.OrderResponse, error)
	CreateManualOrder(ctx context.Context, req *dtos.CreateOrderRequest) (*dtos.OrderResponse, error)
}

type IOrderUpdateUseCase interface {
	UpdateOrder(ctx context.Context, existingOrder *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) (*dtos.OrderResponse, error)
}

type IRequestConfirmationUseCase interface {
	RequestConfirmation(ctx context.Context, orderID string) error
}

type ISendGuideNotificationUseCase interface {
	SendGuideNotification(ctx context.Context, orderID string) error
}

type IOrderStatusUseCase interface {
	ChangeStatus(ctx context.Context, orderID string, req *dtos.ChangeStatusRequest) (*dtos.OrderResponse, error)
}

type IOrderUseCase interface {
	GetOrderByID(ctx context.Context, id string) (*dtos.OrderResponse, error)
	GetOrderRaw(ctx context.Context, id string) (*dtos.OrderRawResponse, error)
	GetOrderHistory(ctx context.Context, orderID string) ([]dtos.OrderHistoryResponse, error)
	ListOrders(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*dtos.OrdersListResponse, error)
	UpdateOrder(ctx context.Context, id string, req *dtos.UpdateOrderRequest) (*dtos.OrderResponse, error)
	DeleteOrder(ctx context.Context, id string) error
}

type IIntegrationEventPublisher interface {
	PublishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
	PublishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
	PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{})
}

type IOrderRabbitPublisher interface {
	PublishOrderCreated(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderUpdated(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderCancelled(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderStatusChanged(ctx context.Context, order *entities.ProbabilityOrder, previousStatus, currentStatus string) error
	PublishConfirmationRequested(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishGuideNotificationRequested(ctx context.Context, order *entities.ProbabilityOrder) error
	PublishOrderEvent(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error
}
