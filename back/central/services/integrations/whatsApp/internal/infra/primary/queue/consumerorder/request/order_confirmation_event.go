package request

// OrderConfirmationEvent representa el evento de confirmación de orden
type OrderConfirmationEvent struct {
	EventType       string  `json:"event_type"`
	OrderID         string  `json:"order_id"`
	OrderNumber     string  `json:"order_number"`
	BusinessID      *uint   `json:"business_id"`
	CustomerName    string  `json:"customer_name"`
	CustomerPhone   string  `json:"customer_phone"`
	CustomerEmail   string  `json:"customer_email"`
	TotalAmount     float64 `json:"total_amount"`
	Currency        string  `json:"currency"`
	ItemsSummary    string  `json:"items_summary"`
	ShippingAddress string  `json:"shipping_address"`
	PaymentMethodID uint    `json:"payment_method_id"`
	IntegrationID   uint    `json:"integration_id"`
	Platform        string  `json:"platform"`
	Timestamp       int64   `json:"timestamp"`
	// Nuevos campos para configuración dinámica
	TemplateName  string `json:"template_name"`   // Nombre de la plantilla a usar
	Language      string `json:"language"`        // Idioma de la plantilla
	RecipientType string `json:"recipient_type"`  // Tipo de destinatario
}
