package responses

import (
	"time"
)

// SSEConnectionResponse representa la respuesta de conexión SSE
type SSEConnectionResponse struct {
	Status        string    `json:"status"`
	IntegrationID int64     `json:"integration_id"`
	BusinessID    string    `json:"business_id"`
	ConnectedAt   time.Time `json:"connected_at"`
	ConnectionID  string    `json:"connection_id"`
}

// SSEStatusResponse representa el estado de las conexiones SSE
type SSEStatusResponse struct {
	IntegrationID     int64     `json:"integration_id"`
	ActiveConnections int       `json:"active_connections"`
	LastActivity      time.Time `json:"last_activity"`
	Status            string    `json:"status"`
}
