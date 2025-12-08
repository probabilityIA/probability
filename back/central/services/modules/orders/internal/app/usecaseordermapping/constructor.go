package usecaseordermapping

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IOrderMappingUseCase interface {
	MapAndSaveOrder(ctx context.Context, dto *domain.CanonicalOrderDTO) (*domain.OrderResponse, error)
}

type UseCaseOrderMapping struct {
	repo           domain.IRepository
	logger         log.ILogger
	eventPublisher domain.IOrderEventPublisher
}

func New(repo domain.IRepository, logger log.ILogger, eventPublisher domain.IOrderEventPublisher) IOrderMappingUseCase {
	return &UseCaseOrderMapping{
		repo:           repo,
		logger:         logger,
		eventPublisher: eventPublisher,
	}
}
