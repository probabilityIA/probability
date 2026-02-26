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
