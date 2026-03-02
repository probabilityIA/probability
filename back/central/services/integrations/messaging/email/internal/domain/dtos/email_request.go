package dtos

// SendEmailDTO contiene los datos necesarios para enviar un email de notificación
type SendEmailDTO struct {
	EventType     string
	BusinessID    uint
	IntegrationID uint
	ConfigID      uint
	CustomerEmail string
	EventData     map[string]interface{}
}
