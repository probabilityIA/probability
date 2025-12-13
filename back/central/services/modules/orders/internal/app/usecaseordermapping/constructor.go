package usecaseordermapping

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UseCaseOrderMapping struct {
	repo           domain.IRepository
	logger         log.ILogger
	eventPublisher domain.IOrderEventPublisher
}

func New(repo domain.IRepository, logger log.ILogger, eventPublisher domain.IOrderEventPublisher) domain.IOrderMappingUseCase {
	return &UseCaseOrderMapping{
		repo:           repo,
		logger:         logger,
		eventPublisher: eventPublisher,
	}
}
