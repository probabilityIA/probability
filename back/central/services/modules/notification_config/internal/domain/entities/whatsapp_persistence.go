package entities

import "time"

type WhatsAppConversation struct {
	ID               string
	PhoneNumber      string
	OrderNumber      string
	ConversationType string
	BusinessID       uint
	CurrentState     string
	LastMessageID    string
	LastTemplateID   string
	Metadata         map[string]interface{}
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ExpiresAt        time.Time
}

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
