package responses

import (
	"time"
)

// PublishEventResponse representa la respuesta de publicación de evento
type PublishEventResponse struct {
	Status      string                 `json:"status"`
	EventID     string                 `json:"event_id"`
	PublishedAt time.Time              `json:"published_at"`
	Recipients  int                    `json:"recipients"`
	Errors      []string               `json:"errors,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ConnectionResponse representa la respuesta de operaciones de conexión
type ConnectionResponse struct {
	Status        string                 `json:"status"`
	IntegrationID int64                  `json:"integration_id"`
	ConnectionID  string                 `json:"connection_id,omitempty"`
	Operation     string                 `json:"operation"`
	Timestamp     time.Time              `json:"timestamp"`
	Message       string                 `json:"message,omitempty"`
	Errors        []string               `json:"errors,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ConnectionInfoResponse representa la información de conexiones
type ConnectionInfoResponse struct {
	IntegrationID     int64                  `json:"integration_id"`
	ActiveConnections int                    `json:"active_connections"`
	TotalConnections  int                    `json:"total_connections"`
	ConnectedAt       time.Time              `json:"connected_at"`
	LastActivity      time.Time              `json:"last_activity"`
	Status            string                 `json:"status"`
	ConnectionIDs     []string               `json:"connection_ids,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}
