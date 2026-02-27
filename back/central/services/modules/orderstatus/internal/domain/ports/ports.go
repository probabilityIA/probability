package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

// IRepository define la interfaz para el almacenamiento de mapeos de estado
// PURO - Solo tipos de dominio, sin tipos de infraestructura
type IRepository interface {
	Create(ctx context.Context, mapping *entities.OrderStatusMapping) error
	GetByID(ctx context.Context, id uint) (*entities.OrderStatusMapping, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.OrderStatusMapping, int64, error)
	Update(ctx context.Context, mapping *entities.OrderStatusMapping) error
	Delete(ctx context.Context, id uint) error
	ToggleActive(ctx context.Context, id uint) (*entities.OrderStatusMapping, error)
	Exists(ctx context.Context, integrationTypeID uint, originalStatus string) (bool, error)
	GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error)
	ListOrderStatuses(ctx context.Context, isActive *bool) ([]entities.OrderStatusInfo, error)
	ListFulfillmentStatuses(ctx context.Context, isActive *bool) ([]entities.FulfillmentStatusInfo, error)

	// CRUD para estados de Probability
	CreateOrderStatus(ctx context.Context, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error)
	GetOrderStatusByID(ctx context.Context, id uint) (*entities.OrderStatusInfo, error)
	UpdateOrderStatus(ctx context.Context, id uint, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error)
	DeleteOrderStatus(ctx context.Context, id uint) error

	// Estados por canal de integración (ecommerce)
	ListEcommerceIntegrationTypes(ctx context.Context, businessID uint) ([]entities.IntegrationTypeInfo, error)
	ListChannelStatuses(ctx context.Context, integrationTypeID uint, isActive *bool) ([]entities.ChannelStatusInfo, error)
	CreateChannelStatus(ctx context.Context, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error)
	GetChannelStatusByID(ctx context.Context, id uint) (*entities.ChannelStatusInfo, error)
	UpdateChannelStatus(ctx context.Context, id uint, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error)
	DeleteChannelStatus(ctx context.Context, id uint) error
}
