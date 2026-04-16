package entities

import "time"

// DeliveryResult representa el resultado del envío de una notificación.
// Se publica a RabbitMQ para que notification_config lo persista en email_logs.
type DeliveryResult struct {
	Channel       string // "email"
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
