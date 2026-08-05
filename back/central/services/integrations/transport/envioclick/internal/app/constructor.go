package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	Quote(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest, metas *[]domain.SyncMeta) (*domain.QuoteResponse, error)
	Generate(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest, metas *[]domain.SyncMeta) (*domain.GenerateResponse, error)
	ResolveCODValue(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest, carrier string, netTarget float64, metas *[]domain.SyncMeta) (float64, bool)
	Track(ctx context.Context, baseURL, apiKey string, trackingNumber string, metas *[]domain.SyncMeta) (*domain.TrackingResponse, error)
	TrackByOrdersBatch(ctx context.Context, baseURL, apiKey string, orders []int64, metas *[]domain.SyncMeta) (*domain.TrackingResponse, error)
	Cancel(ctx context.Context, baseURL, apiKey string, trackingNumber string, idOrder int64, metas *[]domain.SyncMeta) (*domain.CancelResponse, error)
	CancelBatch(ctx context.Context, baseURL, apiKey string, req domain.CancelBatchRequest, metas *[]domain.SyncMeta) (*domain.CancelBatchResponse, error)
}

type useCase struct {
	client domain.IEnvioClickClient
	log    log.ILogger
}

func New(client domain.IEnvioClickClient, logger log.ILogger) IUseCase {
	return &useCase{
		client: client,
		log:    logger.WithModule("transport.envioclick.usecase"),
	}
}
