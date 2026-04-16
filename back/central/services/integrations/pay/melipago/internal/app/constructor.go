package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/pay/melipago/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase interfaz del use case de MercadoPago
type IUseCase interface {
	ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error
}

// useCase implementa IUseCase
type useCase struct {
	meliPagoClient    ports.IMeliPagoClient
	integrationRepo   ports.IIntegrationRepository
	responsePublisher ports.IResponsePublisher
	log               log.ILogger
}

// New crea una nueva instancia del use case de MercadoPago
func New(
	meliPagoClient ports.IMeliPagoClient,
	integrationRepo ports.IIntegrationRepository,
	responsePublisher ports.IResponsePublisher,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		meliPagoClient:    meliPagoClient,
		integrationRepo:   integrationRepo,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("melipago.usecase"),
	}
}
