package request

import "time"

// MessageLogEvent es el payload recibido de RabbitMQ para eventos de message log
type MessageLogEvent struct {
	EventType string `json:"event_type"`
	// Para messagelog.created
	MessageLog *MessageLogPayload `json:"message_log,omitempty"`
	// Para messagelog.status_updated
	Update *StatusUpdatePayload `json:"update,omitempty"`
	// Timestamp del evento
	Timestamp int64 `json:"timestamp"`
}

// MessageLogPayload contiene los datos del message log
type MessageLogPayload struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	Direction      string     `json:"direction"`
	MessageID      string     `json:"message_id"`
	TemplateName   string     `json:"template_name"`
	Content        string     `json:"content"`
	Status         string     `json:"status"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// StatusUpdatePayload contiene los datos de actualización de estado
type StatusUpdatePayload struct {
	MessageID  string            `json:"message_id"`
	Status     string            `json:"status"`
	Timestamps map[string]string `json:"timestamps"`
}
