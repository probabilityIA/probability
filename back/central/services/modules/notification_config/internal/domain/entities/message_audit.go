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
