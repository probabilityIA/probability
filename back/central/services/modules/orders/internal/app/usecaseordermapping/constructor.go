package usecaseordermapping

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorderscore"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UseCaseOrderMapping struct {
	repo                  ports.IRepository
	logger                log.ILogger
	redisEventPublisher   ports.IOrderEventPublisher   // Redis Pub/Sub
	rabbitEventPublisher  ports.IOrderRabbitPublisher  // RabbitMQ
	scoreUseCase          ports.IOrderScoreUseCase
}

func New(
	repo ports.IRepository,
	logger log.ILogger,
	redisPublisher ports.IOrderEventPublisher,
	rabbitPublisher ports.IOrderRabbitPublisher,
) ports.IOrderMappingUseCase {
	return &UseCaseOrderMapping{
		repo:                 repo,
		logger:               logger,
		redisEventPublisher:  redisPublisher,
		rabbitEventPublisher: rabbitPublisher,
		scoreUseCase:         usecaseorderscore.New(repo),
	}
}
