package usecaseupdatestatus

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// UseCaseUpdateStatus maneja los cambios de estado de órdenes con strategy pattern
type UseCaseUpdateStatus struct {
	repo                 ports.IRepository
	logger               log.ILogger
	rabbitEventPublisher ports.IOrderRabbitPublisher
}

// New crea una nueva instancia del caso de uso de cambio de estado
func New(
	repo ports.IRepository,
	logger log.ILogger,
	rabbitPublisher ports.IOrderRabbitPublisher,
) ports.IOrderStatusUseCase {
	return &UseCaseUpdateStatus{
		repo:                 repo,
		logger:               logger,
		rabbitEventPublisher: rabbitPublisher,
	}
}
