package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/transport/mipaquete/internal/domain"
)

// Quote gets shipping rates from MiPaquete
func (uc *UseCase) Quote(ctx context.Context, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
	uc.log.Warn(ctx).Msg("⚠️ MiPaquete Quote not yet implemented")
	return nil, fmt.Errorf("mipaquete: Quote not yet implemented")
}

// Generate creates a shipment and generates a guide with MiPaquete
func (uc *UseCase) Generate(ctx context.Context, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error) {
	uc.log.Warn(ctx).Msg("⚠️ MiPaquete Generate not yet implemented")
	return nil, fmt.Errorf("mipaquete: Generate not yet implemented")
}

// Track gets tracking data from MiPaquete
func (uc *UseCase) Track(ctx context.Context, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
	uc.log.Warn(ctx).Msg("⚠️ MiPaquete Track not yet implemented")
	return nil, fmt.Errorf("mipaquete: Track not yet implemented")
}

// Cancel cancels a shipment in MiPaquete
func (uc *UseCase) Cancel(ctx context.Context, apiKey string, idShipment string) (*domain.CancelResponse, error) {
	uc.log.Warn(ctx).Msg("⚠️ MiPaquete Cancel not yet implemented")
	return nil, fmt.Errorf("mipaquete: Cancel not yet implemented")
}
