package usecaseorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// UseCaseOrder contiene los casos de uso CRUD básicos de órdenes
type UseCaseOrder struct {
	repo                 ports.IRepository
	rabbitEventPublisher ports.IOrderRabbitPublisher
	logger               log.ILogger
}

// New crea una nueva instancia de UseCaseOrder retornando la interfaz IOrderUseCase
func New(
	repo ports.IRepository,
	rabbitPublisher ports.IOrderRabbitPublisher,
	logger log.ILogger,
) ports.IOrderUseCase {
	return &UseCaseOrder{
		repo:                 repo,
		rabbitEventPublisher: rabbitPublisher,
		logger:               logger,
	}
}
