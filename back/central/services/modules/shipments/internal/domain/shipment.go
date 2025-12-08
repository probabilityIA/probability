package domain

import (
	"time"

	"gorm.io/datatypes"
)

// Shipment representa un envío en el dominio
type Shipment struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID string `json:"order_id"`

	TrackingNumber *string `json:"tracking_number"`
	TrackingURL    *string `json:"tracking_url"`
	Carrier        *string `json:"carrier"`
	CarrierCode    *string `json:"carrier_code"`

	GuideID  *string `json:"guide_id"`
	GuideURL *string `json:"guide_url"`

	Status      string     `json:"status"`
	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`

	ShippingAddressID *uint `json:"shipping_address_id"`

	ShippingCost  *float64 `json:"shipping_cost"`
	InsuranceCost *float64 `json:"insurance_cost"`
	TotalCost     *float64 `json:"total_cost"`

	Weight *float64 `json:"weight"`
	Height *float64 `json:"height"`
	Width  *float64 `json:"width"`
	Length *float64 `json:"length"`

	WarehouseID   *uint  `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
	DriverID      *uint  `json:"driver_id"`
	DriverName    string `json:"driver_name"`
	IsLastMile    bool   `json:"is_last_mile"`

	EstimatedDelivery *time.Time     `json:"estimated_delivery"`
	DeliveryNotes     *string        `json:"delivery_notes"`
	Metadata          datatypes.JSON `json:"metadata"`
}

