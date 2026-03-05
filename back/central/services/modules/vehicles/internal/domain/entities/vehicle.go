package entities

import "time"

type Vehicle struct {
	ID                 uint
	BusinessID         uint
	Type               string
	LicensePlate       string
	Brand              string
	VehicleModel       string
	Year               *int
	Color              string
	Status             string
	WeightCapacityKg   *float64
	VolumeCapacityM3   *float64
	PhotoURL           string
	InsuranceExpiry    *time.Time
	RegistrationExpiry *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
