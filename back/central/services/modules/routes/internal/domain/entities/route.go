package entities

import "time"

type Route struct {
	ID                uint
	BusinessID        uint
	DriverID          *uint
	VehicleID         *uint
	Status            string
	Date              time.Time
	StartTime         *time.Time
	EndTime           *time.Time
	ActualStartTime   *time.Time
	ActualEndTime     *time.Time
	OriginWarehouseID *uint
	OriginAddress     string
	OriginLat         *float64
	OriginLng         *float64
	TotalStops        int
	CompletedStops    int
	FailedStops       int
	TotalDistanceKm   *float64
	TotalDurationMin  *int
	Notes             *string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Denormalized
	DriverName   string
	VehiclePlate string

	// Loaded relationships
	Stops []RouteStop
}
