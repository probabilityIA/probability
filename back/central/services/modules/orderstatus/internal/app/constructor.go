package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase define la interfaz para la lógica de negocio de mapeos de estado
type IUseCase interface {
	CreateOrderStatusMapping(ctx context.Context, mapping *entities.OrderStatusMapping) (*entities.OrderStatusMapping, error)
	GetOrderStatusMapping(ctx context.Context, id uint) (*entities.OrderStatusMapping, error)
	ListOrderStatusMappings(ctx context.Context, filters map[string]interface{}) ([]entities.OrderStatusMapping, int64, error)
	UpdateOrderStatusMapping(ctx context.Context, id uint, mapping *entities.OrderStatusMapping) (*entities.OrderStatusMapping, error)
	DeleteOrderStatusMapping(ctx context.Context, id uint) error
	ToggleOrderStatusMappingActive(ctx context.Context, id uint) (*entities.OrderStatusMapping, error)
	ListOrderStatuses(ctx context.Context, isActive *bool) ([]entities.OrderStatusInfo, error)
	ListFulfillmentStatuses(ctx context.Context, isActive *bool) ([]entities.FulfillmentStatusInfo, error)

	// CRUD para estados de Probability
	CreateOrderStatus(ctx context.Context, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error)
	GetOrderStatus(ctx context.Context, id uint) (*entities.OrderStatusInfo, error)
	UpdateOrderStatus(ctx context.Context, id uint, status *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error)
	DeleteOrderStatus(ctx context.Context, id uint) error

	// Estados por canal de integración (ecommerce)
	ListEcommerceIntegrationTypes(ctx context.Context, businessID uint) ([]entities.IntegrationTypeInfo, error)
	ListChannelStatuses(ctx context.Context, integrationTypeID uint, isActive *bool) ([]entities.ChannelStatusInfo, error)
	CreateChannelStatus(ctx context.Context, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error)
	GetChannelStatus(ctx context.Context, id uint) (*entities.ChannelStatusInfo, error)
	UpdateChannelStatus(ctx context.Context, id uint, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error)
	DeleteChannelStatus(ctx context.Context, id uint) error
}

// useCase implementa IUseCase
type useCase struct {
	repo ports.IRepository
	log  log.ILogger
}

// New crea una nueva instancia del caso de uso
func New(repo ports.IRepository, logger log.ILogger) IUseCase {
	return &useCase{
		repo: repo,
		log:  logger.WithModule("orderstatus-usecase"),
	}
}
