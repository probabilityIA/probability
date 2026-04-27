package request

type OrderConfirmationEvent struct {
	EventType         string  `json:"event_type"`
	OrderID           string  `json:"order_id"`
	OrderNumber       string  `json:"order_number"`
	BusinessID        *uint   `json:"business_id"`
	BusinessName      string  `json:"business_name"`
	CustomerName      string  `json:"customer_name"`
	CustomerPhone     string  `json:"customer_phone"`
	CustomerEmail     string  `json:"customer_email"`
	TotalAmount       float64 `json:"total_amount"`
	Currency          string  `json:"currency"`
	ItemsSummary      string  `json:"items_summary"`
	ShippingAddress   string  `json:"shipping_address"`
	ShippingStreet    string  `json:"shipping_street"`
	ShippingCity      string  `json:"shipping_city"`
	ShippingState     string  `json:"shipping_state"`
	PaymentMethodID   uint    `json:"payment_method_id"`
	PaymentMethodName string  `json:"payment_method_name"`
	TrackingNumber    string  `json:"tracking_number"`
	Carrier           string  `json:"carrier"`
	IntegrationID     uint    `json:"integration_id"`
	Platform          string  `json:"platform"`
	Timestamp         int64   `json:"timestamp"`
	TemplateName      string  `json:"template_name"`
	Language          string  `json:"language"`
	RecipientType     string  `json:"recipient_type"`
}
