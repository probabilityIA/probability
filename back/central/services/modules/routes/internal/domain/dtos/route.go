package dtos

import "time"

type ListRoutesParams struct {
	BusinessID uint
	Status     string
	DriverID   *uint
	DateFrom   *time.Time
	DateTo     *time.Time
	Search     string
	Page       int
	PageSize   int
}

func (p ListRoutesParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type CreateRouteDTO struct {
	BusinessID        uint
	DriverID          *uint
	VehicleID         *uint
	Date              time.Time
	StartTime         *time.Time
	EndTime           *time.Time
	OriginWarehouseID *uint
	OriginAddress     string
	OriginLat         *float64
	OriginLng         *float64
	Notes             *string
	Stops             []CreateRouteStopDTO
}

type CreateRouteStopDTO struct {
	OrderID       *string
	Address       string
	City          string
	Lat           *float64
	Lng           *float64
	CustomerName  string
	CustomerPhone string
	DeliveryNotes *string
}

type UpdateRouteDTO struct {
	ID                uint
	BusinessID        uint
	DriverID          *uint
	VehicleID         *uint
	Date              time.Time
	StartTime         *time.Time
	EndTime           *time.Time
	OriginWarehouseID *uint
	OriginAddress     string
	OriginLat         *float64
	OriginLng         *float64
	Notes             *string
}

type AddStopDTO struct {
	RouteID       uint
	BusinessID    uint
	OrderID       *string
	Address       string
	City          string
	Lat           *float64
	Lng           *float64
	CustomerName  string
	CustomerPhone string
	DeliveryNotes *string
}

type UpdateStopDTO struct {
	ID            uint
	RouteID       uint
	BusinessID    uint
	Address       string
	City          string
	Lat           *float64
	Lng           *float64
	CustomerName  string
	CustomerPhone string
	DeliveryNotes *string
}

type UpdateStopStatusDTO struct {
	ID            uint
	RouteID       uint
	BusinessID    uint
	Status        string
	FailureReason *string
	SignatureURL  string
	PhotoURL      string
}

type ReorderStopsDTO struct {
	RouteID    uint
	BusinessID uint
	StopIDs    []uint
}
