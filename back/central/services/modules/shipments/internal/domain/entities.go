package domain

import (
	"time"

	"gorm.io/datatypes"
)

// ───────────────────────────────────────────
//
//	SHIPMENT DTOs
//
// ───────────────────────────────────────────

// CreateShipmentRequest representa la solicitud para crear un envío
type CreateShipmentRequest struct {
	OrderID string `json:"order_id" binding:"required"`

	TrackingNumber *string `json:"tracking_number" binding:"omitempty,max=128"`
	TrackingURL    *string `json:"tracking_url" binding:"omitempty,max=512"`
	Carrier        *string `json:"carrier" binding:"omitempty,max=128"`
	CarrierCode    *string `json:"carrier_code" binding:"omitempty,max=50"`

	GuideID  *string `json:"guide_id" binding:"omitempty,max=128"`
	GuideURL *string `json:"guide_url" binding:"omitempty,max=512"`

	Status      string     `json:"status" binding:"omitempty,oneof=pending in_transit delivered failed"`
	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`

	ShippingAddressID *uint `json:"shipping_address_id"`

	ShippingCost  *float64 `json:"shipping_cost" binding:"omitempty,min=0"`
	InsuranceCost *float64 `json:"insurance_cost" binding:"omitempty,min=0"`
	TotalCost     *float64 `json:"total_cost" binding:"omitempty,min=0"`

	Weight *float64 `json:"weight" binding:"omitempty,min=0"`
	Height *float64 `json:"height" binding:"omitempty,min=0"`
	Width  *float64 `json:"width" binding:"omitempty,min=0"`
	Length *float64 `json:"length" binding:"omitempty,min=0"`

	WarehouseID   *uint  `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name" binding:"omitempty,max=128"`
	DriverID      *uint  `json:"driver_id"`
	DriverName    string `json:"driver_name" binding:"omitempty,max=255"`
	IsLastMile    bool   `json:"is_last_mile"`

	EstimatedDelivery *time.Time     `json:"estimated_delivery"`
	DeliveryNotes     *string        `json:"delivery_notes"`
	Metadata          datatypes.JSON `json:"metadata"`
}

// UpdateShipmentRequest representa la solicitud para actualizar un envío
type UpdateShipmentRequest struct {
	TrackingNumber *string `json:"tracking_number" binding:"omitempty,max=128"`
	TrackingURL    *string `json:"tracking_url" binding:"omitempty,max=512"`
	Carrier        *string `json:"carrier" binding:"omitempty,max=128"`
	CarrierCode    *string `json:"carrier_code" binding:"omitempty,max=50"`

	GuideID  *string `json:"guide_id" binding:"omitempty,max=128"`
	GuideURL *string `json:"guide_url" binding:"omitempty,max=512"`

	Status      *string    `json:"status" binding:"omitempty,oneof=pending in_transit delivered failed"`
	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`

	ShippingAddressID *uint `json:"shipping_address_id"`

	ShippingCost  *float64 `json:"shipping_cost" binding:"omitempty,min=0"`
	InsuranceCost *float64 `json:"insurance_cost" binding:"omitempty,min=0"`
	TotalCost     *float64 `json:"total_cost" binding:"omitempty,min=0"`

	Weight *float64 `json:"weight" binding:"omitempty,min=0"`
	Height *float64 `json:"height" binding:"omitempty,min=0"`
	Width  *float64 `json:"width" binding:"omitempty,min=0"`
	Length *float64 `json:"length" binding:"omitempty,min=0"`

	WarehouseID   *uint   `json:"warehouse_id"`
	WarehouseName *string `json:"warehouse_name" binding:"omitempty,max=128"`
	DriverID      *uint   `json:"driver_id"`
	DriverName    *string `json:"driver_name" binding:"omitempty,max=255"`
	IsLastMile    *bool   `json:"is_last_mile"`

	EstimatedDelivery *time.Time     `json:"estimated_delivery"`
	DeliveryNotes     *string        `json:"delivery_notes"`
	Metadata          datatypes.JSON `json:"metadata"`
}

// ShipmentResponse representa la respuesta de un envío
type ShipmentResponse struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID string `json:"order_id"`

	TrackingNumber *string `json:"tracking_number,omitempty"`
	TrackingURL    *string `json:"tracking_url,omitempty"`
	Carrier        *string `json:"carrier,omitempty"`
	CarrierCode    *string `json:"carrier_code,omitempty"`

	GuideID  *string `json:"guide_id,omitempty"`
	GuideURL *string `json:"guide_url,omitempty"`

	Status      string     `json:"status"`
	ShippedAt   *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`

	ShippingAddressID *uint `json:"shipping_address_id,omitempty"`

	ShippingCost  *float64 `json:"shipping_cost,omitempty"`
	InsuranceCost *float64 `json:"insurance_cost,omitempty"`
	TotalCost     *float64 `json:"total_cost,omitempty"`

	Weight *float64 `json:"weight,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Width  *float64 `json:"width,omitempty"`
	Length *float64 `json:"length,omitempty"`

	WarehouseID   *uint  `json:"warehouse_id,omitempty"`
	WarehouseName string `json:"warehouse_name"`
	DriverID      *uint  `json:"driver_id,omitempty"`
	DriverName    string `json:"driver_name"`
	IsLastMile    bool   `json:"is_last_mile"`

	EstimatedDelivery *time.Time     `json:"estimated_delivery,omitempty"`
	DeliveryNotes     *string        `json:"delivery_notes,omitempty"`
	Metadata          datatypes.JSON `json:"metadata,omitempty"`
}

// ShipmentsListResponse representa la respuesta paginada de envíos
type ShipmentsListResponse struct {
	Data       []ShipmentResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

