package dtos

// ShipmentSummary representa un resumen del envío para la orden
// ✅ DTO PURO - SIN TAGS
type ShipmentSummary struct {
	ID             uint
	Carrier        *string
	TrackingNumber *string
	GuideURL       *string
	Status         string
	TotalCost      *float64
}
