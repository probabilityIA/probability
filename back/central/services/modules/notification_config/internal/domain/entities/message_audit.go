package entities

import "time"

type MessageAuditLog struct {
	ID             string
	ConversationID string
	MessageID      string
	Direction      string
	TemplateName   string
	Content        string
	Status         string
	DeliveredAt    *time.Time
	ReadAt         *time.Time
	CreatedAt      time.Time

	PhoneNumber string
	OrderNumber string
	BusinessID  uint
}

type MessageAuditStats struct {
	TotalSent      int64
	TotalDelivered int64
	TotalRead      int64
	TotalFailed    int64
	SuccessRate    float64
}

type ConversationSummary struct {
	ID                   string
	PhoneNumber          string
	OrderNumber          string
	ConversationType     string
	BusinessID           uint
	CurrentState         string
	MessageCount         int
	LastMessageContent   string
	LastMessageDirection string
	LastMessageStatus    string
	LastActivity         time.Time
	CreatedAt            time.Time
}

type ConversationMessage struct {
	ID           string
	Direction    string
	MessageID    string
	TemplateName string
	Content      string
	Status       string
	DeliveredAt  *time.Time
	ReadAt       *time.Time
	CreatedAt    time.Time
}
