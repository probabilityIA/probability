package request

import "time"

// ConversationEvent es el payload recibido de RabbitMQ para eventos de conversación
type ConversationEvent struct {
	EventType    string               `json:"event_type"`
	Conversation ConversationPayload  `json:"conversation"`
	Timestamp    int64                `json:"timestamp"`
}

// ConversationPayload contiene los datos de la conversación
type ConversationPayload struct {
	ID             string                 `json:"id"`
	PhoneNumber    string                 `json:"phone_number"`
	OrderNumber    string                 `json:"order_number"`
	BusinessID     uint                   `json:"business_id"`
	CurrentState   string                 `json:"current_state"`
	LastMessageID  string                 `json:"last_message_id"`
	LastTemplateID string                 `json:"last_template_id"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ExpiresAt      time.Time              `json:"expires_at"`
}
