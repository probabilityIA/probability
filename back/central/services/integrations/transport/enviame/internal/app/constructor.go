package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase defines the operations available for the Enviame transport provider.
type IUseCase interface {
	Quote(ctx context.Context, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error)
	Generate(ctx context.Context, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error)
	Track(ctx context.Context, apiKey string, trackingNumber string) (*domain.TrackingResponse, error)
	Cancel(ctx context.Context, apiKey string, idShipment string) (*domain.CancelResponse, error)
}

// useCase implements the transport operations for Enviame
type useCase struct {
	client domain.IEnviameClient
	log    log.ILogger
}

// New creates the Enviame transport use case
func New(
	client domain.IEnviameClient,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		client: client,
		log:    logger.WithModule("enviame.usecase"),
	}
}
