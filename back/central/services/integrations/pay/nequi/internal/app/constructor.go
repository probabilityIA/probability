package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase interfaz del use case de Nequi
type IUseCase interface {
	ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error
}

// useCase implementa IUseCase
type useCase struct {
	nequiClient       ports.INequiClient
	integrationRepo   ports.IIntegrationRepository
	responsePublisher ports.IResponsePublisher
	log               log.ILogger
}

// New crea una nueva instancia del use case de Nequi
func New(
	nequiClient ports.INequiClient,
	integrationRepo ports.IIntegrationRepository,
	responsePublisher ports.IResponsePublisher,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		nequiClient:       nequiClient,
		integrationRepo:   integrationRepo,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("nequi.usecase"),
	}
}
