package request

type PushEvent struct {
	EventType     string  `json:"event_type"`
	EventCategory string  `json:"event_category"`
	BusinessID    uint    `json:"business_id"`
	IntegrationID uint    `json:"integration_id"`
	ConfigID      uint    `json:"config_id"`
	OrderID       string  `json:"order_id"`
	OrderNumber   string  `json:"order_number"`
	ShipmentID    string  `json:"shipment_id"`
	TrackingNo    string  `json:"tracking_number"`
	Carrier       string  `json:"carrier"`
	CustomerName  string  `json:"customer_name"`
	StatusCode    string  `json:"status_code"`
	StatusName    string  `json:"status_name"`
	CodTotal      float64 `json:"cod_total"`
	Balance       float64 `json:"balance"`
}
