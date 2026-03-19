package request

// ShipmentGuideEvent representa el evento de guía de envío generada para WhatsApp
type ShipmentGuideEvent struct {
	EventType      string `json:"event_type"`
	BusinessID     *uint  `json:"business_id"`
	IntegrationID  uint   `json:"integration_id"`
	ConfigID       uint   `json:"config_id"`
	ShipmentID     uint   `json:"shipment_id"`
	TrackingNumber string `json:"tracking_number"`
	LabelURL       string `json:"label_url"`
	Carrier        string `json:"carrier"`
	CustomerName   string `json:"customer_name"`
	CustomerPhone  string `json:"customer_phone"`
	OrderNumber    string `json:"order_number"`
	BusinessName   string `json:"business_name"`
	CorrelationID  string `json:"correlation_id"`
}
