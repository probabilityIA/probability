package request

// EmailNotificationEvent es el formato JSON que se recibe de la cola RabbitMQ
type EmailNotificationEvent struct {
	EventType     string                 `json:"event_type"`
	BusinessID    uint                   `json:"business_id"`
	IntegrationID uint                   `json:"integration_id"`
	ConfigID      uint                   `json:"config_id"`
	CustomerEmail string                 `json:"customer_email"`
	EventData     map[string]interface{} `json:"event_data"`
}
