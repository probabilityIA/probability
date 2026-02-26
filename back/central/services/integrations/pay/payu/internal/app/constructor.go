package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/pay/payu/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase interfaz del use case de PayU
type IUseCase interface {
	ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error
}

// useCase implementa IUseCase
type useCase struct {
	payuClient        ports.IPayUClient
	integrationRepo   ports.IIntegrationRepository
	responsePublisher ports.IResponsePublisher
	log               log.ILogger
}

// New crea una nueva instancia del use case de PayU
func New(
	payuClient ports.IPayUClient,
	integrationRepo ports.IIntegrationRepository,
	responsePublisher ports.IResponsePublisher,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		payuClient:        payuClient,
		integrationRepo:   integrationRepo,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("payu.usecase"),
	}
}
