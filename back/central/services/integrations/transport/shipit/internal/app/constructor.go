package app

import (
	"context"
	"sync"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	Quote(ctx context.Context, baseURL string, creds domain.Credentials, req domain.GuideRequest) (*domain.QuoteResponse, error)
	Generate(ctx context.Context, baseURL string, creds domain.Credentials, req domain.GuideRequest, sandbox bool) (*domain.GenerateResponse, error)
	Track(ctx context.Context, baseURL string, creds domain.Credentials, trackingNumber string) (*domain.TrackOutput, error)
	TestAuthentication(ctx context.Context, baseURL string, creds domain.Credentials) error
}

type useCase struct {
	client        domain.IShipitClient
	log           log.ILogger
	communesMu    sync.Mutex
	communesCache map[string]map[string]uint
}

func New(client domain.IShipitClient, logger log.ILogger) IUseCase {
	return &useCase{
		client:        client,
		log:           logger.WithModule("transport.shipit.usecase"),
		communesCache: make(map[string]map[string]uint),
	}
}
