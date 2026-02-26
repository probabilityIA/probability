package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/pay/stripe/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase interfaz del use case de Stripe
type IUseCase interface {
	ProcessPayment(ctx context.Context, msg *PaymentRequestMsg) error
}

// useCase implementa IUseCase
type useCase struct {
	stripeClient      ports.IStripeClient
	integrationRepo   ports.IIntegrationRepository
	responsePublisher ports.IResponsePublisher
	log               log.ILogger
}

// New crea una nueva instancia del use case de Stripe
func New(
	stripeClient ports.IStripeClient,
	integrationRepo ports.IIntegrationRepository,
	responsePublisher ports.IResponsePublisher,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		stripeClient:      stripeClient,
		integrationRepo:   integrationRepo,
		responsePublisher: responsePublisher,
		log:               logger.WithModule("stripe.usecase"),
	}
}
