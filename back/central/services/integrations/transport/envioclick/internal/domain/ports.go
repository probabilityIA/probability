package domain

// IEnvioClickClient defines the contract for the EnvioClick HTTP client.
// baseURL is the production URL from integration_types.base_url (falls back to DefaultBaseURL if empty).
// apiKey supports per-business credentials or the platform shared key.
type IEnvioClickClient interface {
	Quote(baseURL, apiKey string, req QuoteRequest) (*QuoteResponse, error)
	Generate(baseURL, apiKey string, req QuoteRequest) (*GenerateResponse, error)
	Track(baseURL, apiKey string, trackingNumber string) (*TrackingResponse, error)
	Cancel(baseURL, apiKey string, idShipment string) (*CancelResponse, error)
}
