package entities

import "time"

// WhatsAppConversation es la entidad de dominio para persistir conversaciones de WhatsApp.
// Los datos provienen de eventos RabbitMQ publicados por el módulo WhatsApp.
type WhatsAppConversation struct {
	ID             string
	PhoneNumber    string
	OrderNumber    string
	BusinessID     uint
	CurrentState   string
	LastMessageID  string
	LastTemplateID string
	Metadata       map[string]interface{}
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExpiresAt      time.Time
}

// WhatsAppMessageLogEntry es la entidad de dominio para persistir message logs de WhatsApp.
type WhatsAppMessageLogEntry struct {
	ID             string
	ConversationID string
	Direction      string
	MessageID      string
	TemplateName   string
	Content        string
	Status         string
	DeliveredAt    *time.Time
	ReadAt         *time.Time
	CreatedAt      time.Time
}
