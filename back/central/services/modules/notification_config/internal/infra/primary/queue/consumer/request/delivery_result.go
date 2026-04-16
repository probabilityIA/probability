package request

import "time"

// DeliveryResult representa el mensaje JSON que publica el módulo de email
// después de enviar (o fallar) una notificación.
type DeliveryResult struct {
	Channel       string    `json:"Channel"`
	BusinessID    uint      `json:"BusinessID"`
	IntegrationID uint      `json:"IntegrationID"`
	ConfigID      uint      `json:"ConfigID"`
	To            string    `json:"To"`
	Subject       string    `json:"Subject"`
	EventType     string    `json:"EventType"`
	Status        string    `json:"Status"`
	ErrorMessage  string    `json:"ErrorMessage"`
	SentAt        time.Time `json:"SentAt"`
}
