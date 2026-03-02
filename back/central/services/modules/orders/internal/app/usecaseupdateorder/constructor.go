package usecaseupdateorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorderscore"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UseCaseUpdateOrder struct {
	repo                      ports.IRepository
	logger                    log.ILogger
	rabbitEventPublisher      ports.IOrderRabbitPublisher
	integrationEventPublisher ports.IIntegrationEventPublisher
	scoreUseCase              ports.IOrderScoreUseCase
}

func New(
	repo ports.IRepository,
	logger log.ILogger,
	rabbitPublisher ports.IOrderRabbitPublisher,
	integrationEventPub ports.IIntegrationEventPublisher,
) ports.IOrderUpdateUseCase {
	return &UseCaseUpdateOrder{
		repo:                      repo,
		logger:                    logger,
		rabbitEventPublisher:      rabbitPublisher,
		integrationEventPublisher: integrationEventPub,
		scoreUseCase:              usecaseorderscore.New(repo),
	}
}
