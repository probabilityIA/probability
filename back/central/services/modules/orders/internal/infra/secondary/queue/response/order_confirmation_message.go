package response

// OrderConfirmationMessage estructura tipada para mensajes de confirmación WhatsApp
// Sincronizada con whatsapp/internal/infra/primary/queue/consumerorder/request/order_confirmation_event.go
type OrderConfirmationMessage struct {
	EventType        string  `json:"event_type"`
	OrderID          string  `json:"order_id"`
	OrderNumber      string  `json:"order_number"`
	BusinessID       *uint   `json:"business_id"`
	CustomerName     string  `json:"customer_name"`
	CustomerPhone    string  `json:"customer_phone"`
	CustomerEmail    string  `json:"customer_email,omitempty"`
	TotalAmount      float64 `json:"total_amount"`
	Currency         string  `json:"currency"`
	ItemsSummary     string  `json:"items_summary"`
	ShippingAddress  string  `json:"shipping_address"`
	PaymentMethodID  uint    `json:"payment_method_id"`
	IntegrationID    uint    `json:"integration_id"`
	Platform         string  `json:"platform"`
	Timestamp        int64   `json:"timestamp"`
	TemplateName     string  `json:"template_name,omitempty"`
	Language         string  `json:"language,omitempty"`
	RecipientType    string  `json:"recipient_type,omitempty"`
}
