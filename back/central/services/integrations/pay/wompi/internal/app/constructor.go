package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/pay/wompi/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase interfaz del use case de Wompi
type IUseCase interface {
	ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error
}

// useCase implementa IUseCase
type useCase struct {
	wompiClient       ports.IWompiClient
	integrationRepo   ports.IIntegrationRepository
	responsePublisher ports.IResponsePublisher
	log               log.ILogger
}

// New crea una nueva instancia del use case de Wompi
func New(
	wompiClient ports.IWompiClient,
	integrationRepo ports.IIntegrationRepository,
	responsePublisher ports.IResponsePublisher,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		wompiClient:       wompiClient,
		integrationRepo:   integrationRepo,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("wompi.usecase"),
	}
}
