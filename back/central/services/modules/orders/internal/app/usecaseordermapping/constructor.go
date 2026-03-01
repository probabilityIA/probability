package usecaseordermapping

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorderscore"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UseCaseOrderMapping struct {
	repo                     ports.IRepository
	logger                   log.ILogger
	redisEventPublisher      ports.IOrderEventPublisher
	rabbitEventPublisher     ports.IOrderRabbitPublisher
	integrationEventPublisher ports.IIntegrationEventPublisher
	scoreUseCase             ports.IOrderScoreUseCase
}

func New(
	repo ports.IRepository,
	logger log.ILogger,
	redisPublisher ports.IOrderEventPublisher,
	rabbitPublisher ports.IOrderRabbitPublisher,
	integrationEventPub ports.IIntegrationEventPublisher,
) ports.IOrderMappingUseCase {
	return &UseCaseOrderMapping{
		repo:                      repo,
		logger:                    logger,
		redisEventPublisher:       redisPublisher,
		rabbitEventPublisher:      rabbitPublisher,
		integrationEventPublisher: integrationEventPub,
		scoreUseCase:              usecaseorderscore.New(repo),
	}
}
