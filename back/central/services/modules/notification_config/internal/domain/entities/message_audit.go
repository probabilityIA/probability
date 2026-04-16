package entities

import "time"

// MessageAuditLog representa un registro de auditoría de mensajes enviados
// Enriquecido con datos de la conversación (JOIN)
type MessageAuditLog struct {
	ID             string
	ConversationID string
	MessageID      string
	Direction      string // "outbound" | "inbound"
	TemplateName   string
	Content        string
	Status         string // "sent" | "delivered" | "read" | "failed"
	DeliveredAt    *time.Time
	ReadAt         *time.Time
	CreatedAt      time.Time

	// Enriched from conversation JOIN
	PhoneNumber string
	OrderNumber string
	BusinessID  uint
}

// MessageAuditStats contiene estadísticas agregadas de mensajes
type MessageAuditStats struct {
	TotalSent      int64
	TotalDelivered int64
	TotalRead      int64
	TotalFailed    int64
	SuccessRate    float64
}

// ConversationSummary representa el resumen de una conversación para la vista de lista
type ConversationSummary struct {
	ID                   string
	PhoneNumber          string
	OrderNumber          string
	BusinessID           uint
	CurrentState         string
	MessageCount         int
	LastMessageContent   string
	LastMessageDirection string
	LastMessageStatus    string
	LastActivity         time.Time
	CreatedAt            time.Time
}

// ConversationMessage representa un mensaje dentro de una conversación para la vista de chat
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
