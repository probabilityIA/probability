package request

import "time"

type PersistenceEvent struct {
	EventType    string               `json:"event_type"`
	Conversation *ConversationPayload `json:"conversation,omitempty"`
	MessageLog   *MessageLogPayload   `json:"message_log,omitempty"`
	Update       *StatusUpdatePayload `json:"update,omitempty"`
	Timestamp    int64                `json:"timestamp"`
}

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

type StatusUpdatePayload struct {
	MessageID  string            `json:"message_id"`
	Status     string            `json:"status"`
	Timestamps map[string]string `json:"timestamps"`
}
