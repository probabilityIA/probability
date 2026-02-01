package usecaseorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// UseCaseOrder contiene los casos de uso CRUD básicos de órdenes
type UseCaseOrder struct {
	repo                  ports.IRepository
	redisEventPublisher   ports.IOrderEventPublisher   // Redis Pub/Sub
	rabbitEventPublisher  ports.IOrderRabbitPublisher  // RabbitMQ
	logger                log.ILogger
	scoreUseCase          ports.IOrderScoreUseCase
}

// New crea una nueva instancia de UseCaseOrder retornando la interfaz IOrderUseCase
func New(
	repo ports.IRepository,
	redisPublisher ports.IOrderEventPublisher,
	rabbitPublisher ports.IOrderRabbitPublisher,
	logger log.ILogger,
	scoreUseCase ports.IOrderScoreUseCase,
) ports.IOrderUseCase {
	return &UseCaseOrder{
		repo:                 repo,
		redisEventPublisher:  redisPublisher,
		rabbitEventPublisher: rabbitPublisher,
		logger:               logger,
		scoreUseCase:         scoreUseCase,
	}
}
