package entities

import "time"

// ProbabilityShipment representa un envío que se guarda en la base de datos
// ✅ ENTIDAD PURA - SIN TAGS
type ProbabilityShipment struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	OrderID *string

	TrackingNumber *string
	TrackingURL    *string
	Carrier        *string
	CarrierCode    *string

	GuideID  *string
	GuideURL *string

	Status      string
	ShippedAt   *time.Time
	DeliveredAt *time.Time

	ShippingAddressID *uint

	ShippingCost  *float64
	InsuranceCost *float64
	TotalCost     *float64

	Weight *float64
	Height *float64
	Width  *float64
	Length *float64

	WarehouseID   *uint
	WarehouseName string
	DriverID      *uint
	DriverName    string
	IsLastMile    bool

	EstimatedDelivery *time.Time
	DeliveryNotes     *string
	Metadata          []byte
}
