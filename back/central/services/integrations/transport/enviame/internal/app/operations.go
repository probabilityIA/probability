package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/domain"
)

// Quote gets shipping rates from Enviame
func (uc *useCase) Quote(ctx context.Context, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
	uc.log.Warn(ctx).Msg("⚠️ Enviame Quote not yet implemented")
	return nil, fmt.Errorf("enviame: Quote not yet implemented")
}

// Generate creates a shipment and generates a guide with Enviame
func (uc *useCase) Generate(ctx context.Context, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error) {
	uc.log.Warn(ctx).Msg("⚠️ Enviame Generate not yet implemented")
	return nil, fmt.Errorf("enviame: Generate not yet implemented")
}

// Track gets tracking data from Enviame
func (uc *useCase) Track(ctx context.Context, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
	uc.log.Warn(ctx).Msg("⚠️ Enviame Track not yet implemented")
	return nil, fmt.Errorf("enviame: Track not yet implemented")
}

// Cancel cancels a shipment in Enviame
func (uc *useCase) Cancel(ctx context.Context, apiKey string, idShipment string) (*domain.CancelResponse, error) {
	uc.log.Warn(ctx).Msg("⚠️ Enviame Cancel not yet implemented")
	return nil, fmt.Errorf("enviame: Cancel not yet implemented")
}
