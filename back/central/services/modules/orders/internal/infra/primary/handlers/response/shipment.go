package response

// ShipmentSummary representa un resumen del envío
// ✅ DTO HTTP - CON TAGS (json)
type ShipmentSummary struct {
	ID             uint    `json:"id"`
	Carrier        *string `json:"carrier,omitempty"`
	TrackingNumber *string `json:"tracking_number,omitempty"`
	GuideURL       *string `json:"guide_url,omitempty"`
	Status         string  `json:"status"`
}
