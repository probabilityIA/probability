package domain

import "context"

// IEnviameClient defines the contract for the Enviame HTTP client.
// Each method receives an apiKey to support per-business credentials.
type IEnviameClient interface {
	Quote(apiKey string, req QuoteRequest) (*QuoteResponse, error)
	Generate(apiKey string, req QuoteRequest) (*GenerateResponse, error)
	Track(apiKey string, trackingNumber string) (*TrackingResponse, error)
	Cancel(apiKey string, idShipment string) (*CancelResponse, error)
	TestAuthentication(ctx context.Context, apiKey string) error
}
