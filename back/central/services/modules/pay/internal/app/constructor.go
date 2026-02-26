package app

import (
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// useCase implementa todos los casos de uso del módulo de pagos
type useCase struct {
	repo             ports.IRepository
	requestPublisher ports.IRequestPublisher
	ssePublisher     ports.ISSEPublisher
	log              log.ILogger
}

// New crea una nueva instancia del use case de pagos
func New(
	repo ports.IRepository,
	requestPublisher ports.IRequestPublisher,
	ssePublisher ports.ISSEPublisher,
	logger log.ILogger,
) ports.IUseCase {
	return &useCase{
		repo:             repo,
		requestPublisher: requestPublisher,
		ssePublisher:     ssePublisher,
		log:              logger.WithModule("pay.usecase"),
	}
}
