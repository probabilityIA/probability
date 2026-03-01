package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IUseCase defines the operations available for the EnvioClick transport provider.
type IUseCase interface {
	Quote(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error)
	Generate(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error)
	Track(ctx context.Context, baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error)
	Cancel(ctx context.Context, baseURL, apiKey string, idShipment string) (*domain.CancelResponse, error)
}

// useCase handles EnvioClick transport operations
type useCase struct {
	client domain.IEnvioClickClient
	log    log.ILogger
}

// New creates a new EnvioClick use case
func New(client domain.IEnvioClickClient, logger log.ILogger) IUseCase {
	return &useCase{
		client: client,
		log:    logger.WithModule("transport.envioclick.usecase"),
	}
}
