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

// DriverOption is a simplified driver for selection dropdowns
type DriverOption struct {
	ID             uint
	FirstName      string
	LastName       string
	Phone          string
	Identification string
	Status         string
	LicenseType    string
}

// VehicleOption is a simplified vehicle for selection dropdowns
type VehicleOption struct {
	ID           uint
	Type         string
	LicensePlate string
	Brand        string
	VehicleModel string
	Status       string
}

// AssignableOrder represents an order available for route assignment
type AssignableOrder struct {
	ID            string
	OrderNumber   string
	CustomerName  string
	CustomerPhone string
	Address       string
	City          string
	Lat           *float64
	Lng           *float64
	TotalAmount   float64
	ItemCount     int
	CreatedAt     time.Time
}
