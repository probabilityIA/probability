package request

import "time"

type CreateRouteRequest struct {
	DriverID          *uint              `json:"driver_id"`
	VehicleID         *uint              `json:"vehicle_id"`
	Date              time.Time          `json:"date" binding:"required"`
	StartTime         *time.Time         `json:"start_time"`
	EndTime           *time.Time         `json:"end_time"`
	OriginWarehouseID *uint              `json:"origin_warehouse_id"`
	OriginAddress     string             `json:"origin_address" binding:"omitempty,max=500"`
	OriginLat         *float64           `json:"origin_lat"`
	OriginLng         *float64           `json:"origin_lng"`
	Notes             *string            `json:"notes"`
	Stops             []CreateStopRequest `json:"stops"`
}

type CreateStopRequest struct {
	OrderID       *string `json:"order_id"`
	Address       string  `json:"address" binding:"required,max=500"`
	City          string  `json:"city" binding:"omitempty,max=128"`
	Lat           *float64 `json:"lat"`
	Lng           *float64 `json:"lng"`
	CustomerName  string  `json:"customer_name" binding:"omitempty,max=255"`
	CustomerPhone string  `json:"customer_phone" binding:"omitempty,max=50"`
	DeliveryNotes *string `json:"delivery_notes"`
}

type UpdateRouteRequest struct {
	DriverID          *uint      `json:"driver_id"`
	VehicleID         *uint      `json:"vehicle_id"`
	Date              time.Time  `json:"date" binding:"required"`
	StartTime         *time.Time `json:"start_time"`
	EndTime           *time.Time `json:"end_time"`
	OriginWarehouseID *uint      `json:"origin_warehouse_id"`
	OriginAddress     string     `json:"origin_address" binding:"omitempty,max=500"`
	OriginLat         *float64   `json:"origin_lat"`
	OriginLng         *float64   `json:"origin_lng"`
	Notes             *string    `json:"notes"`
}
