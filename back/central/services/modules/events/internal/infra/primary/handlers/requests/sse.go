package requests

// SSEConnectionRequest representa la solicitud para conectar SSE
type SSEConnectionRequest struct {
	IntegrationID int64  `json:"integration_id"`
	BusinessID    string `json:"business_id"`
}
