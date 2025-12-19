package usecaseenvioclick

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/secondary/envioclick"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UseCaseEnvioClick struct {
	logger           log.ILogger
	envioclickClient *envioclick.Client
}

func New(logger log.ILogger, client *envioclick.Client) *UseCaseEnvioClick {
	return &UseCaseEnvioClick{
		logger:           logger,
		envioclickClient: client,
	}
}

func (uc *UseCaseEnvioClick) QuoteShipment(ctx context.Context, req domain.EnvioClickQuoteRequest) (*domain.EnvioClickQuoteResponse, error) {
	uc.logger.Info().Msg("Quoting shipment with EnvioClick")
	// Here we could add logic to validate or enrich the request if needed
	return uc.envioclickClient.Quote(req)
}

func (uc *UseCaseEnvioClick) GenerateGuide(ctx context.Context, req interface{}) (*domain.EnvioClickGenerateResponse, error) {
	uc.logger.Info().Msg("Generating guide with EnvioClick")
	// Here we could add logic to save the guide info to the database (update the order)
	// For now we just pass through to the client
	return uc.envioclickClient.Generate(req)
}
