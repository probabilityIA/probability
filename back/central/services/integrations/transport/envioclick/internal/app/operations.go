package app

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

// Quote gets shipping rates from EnvioClick
func (uc *useCase) Quote(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
	uc.log.Info(ctx).Msg("Quoting shipment with EnvioClick")
	return uc.client.Quote(baseURL, apiKey, req)
}

// Generate creates a shipment and generates a guide with EnvioClick
func (uc *useCase) Generate(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error) {
	uc.log.Info(ctx).Msg("Generating guide with EnvioClick")
	return uc.client.Generate(baseURL, apiKey, req)
}

// Track gets tracking data from EnvioClick
func (uc *useCase) Track(ctx context.Context, baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
	uc.log.Info(ctx).Str("tracking_number", trackingNumber).Msg("Tracking shipment with EnvioClick")
	return uc.client.Track(baseURL, apiKey, trackingNumber)
}

// Cancel cancels a shipment in EnvioClick
func (uc *useCase) Cancel(ctx context.Context, baseURL, apiKey string, idShipment string) (*domain.CancelResponse, error) {
	uc.log.Info(ctx).Str("id_shipment", idShipment).Msg("Canceling shipment with EnvioClick")
	return uc.client.Cancel(baseURL, apiKey, idShipment)
}

// toMap converts a struct to map[string]interface{} via JSON marshaling
func toMap(v interface{}) map[string]interface{} {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}
