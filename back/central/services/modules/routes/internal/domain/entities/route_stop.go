package entities

import "time"

type RouteStop struct {
	ID               uint
	RouteID          uint
	OrderID          *string
	Sequence         int
	Status           string
	Address          string
	City             string
	Lat              *float64
	Lng              *float64
	CustomerName     string
	CustomerPhone    string
	EstimatedArrival *time.Time
	ActualArrival    *time.Time
	ActualDeparture  *time.Time
	SignatureURL     string
	PhotoURL         string
	DeliveryNotes    *string
	FailureReason    *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
