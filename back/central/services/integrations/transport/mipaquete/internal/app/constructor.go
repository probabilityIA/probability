package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/transport/mipaquete/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase defines the operations available for the MiPaquete transport provider.
type IUseCase interface {
	Quote(ctx context.Context, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error)
	Generate(ctx context.Context, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error)
	Track(ctx context.Context, apiKey string, trackingNumber string) (*domain.TrackingResponse, error)
	Cancel(ctx context.Context, apiKey string, idShipment string) (*domain.CancelResponse, error)
}

// useCase implements the transport operations for MiPaquete
type useCase struct {
	client domain.IMiPaqueteClient
	log    log.ILogger
}

// New creates the MiPaquete transport use case
func New(
	client domain.IMiPaqueteClient,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		client: client,
		log:    logger.WithModule("mipaquete.usecase"),
	}
}
