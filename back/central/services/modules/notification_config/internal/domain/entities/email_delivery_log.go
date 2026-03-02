package entities

import "time"

// EmailDeliveryLog representa un log de entrega de email de notificación.
// Sin tags GORM — dominio puro. El mapper convierte a models.EmailLog.
type EmailDeliveryLog struct {
	BusinessID    uint
	IntegrationID uint
	ConfigID      uint
	To            string
	Subject       string
	EventType     string
	Status        string // "sent" | "failed"
	ErrorMessage  string
	SentAt        time.Time
}
