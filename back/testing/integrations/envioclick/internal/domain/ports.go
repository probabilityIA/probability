package domain

// IAPISimulator define la interfaz para el simulador de la API de EnvioClick
type IAPISimulator interface {
	HandleQuote(req QuoteRequest) (*QuoteResponse, error)
	HandleGenerate(req QuoteRequest) (*GenerateResponse, error)
	HandleTrack(trackingNumber string) (*TrackingResponse, error)
	HandleCancel(shipmentID string) (*CancelResponse, error)
}
