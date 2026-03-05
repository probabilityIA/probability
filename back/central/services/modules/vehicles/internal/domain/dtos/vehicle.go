package dtos

import "time"

type ListVehiclesParams struct {
	BusinessID uint
	Search     string
	Type       string
	Status     string
	Page       int
	PageSize   int
}

func (p ListVehiclesParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type CreateVehicleDTO struct {
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
}

type UpdateVehicleDTO struct {
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
}
