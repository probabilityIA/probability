package request

import "time"

type CreateVehicleRequest struct {
	Type               string     `json:"type" binding:"required,oneof=motorcycle car van truck"`
	LicensePlate       string     `json:"license_plate" binding:"required,max=20"`
	Brand              string     `json:"brand" binding:"omitempty,max=64"`
	VehicleModel       string     `json:"vehicle_model" binding:"omitempty,max=64"`
	Year               *int       `json:"year"`
	Color              string     `json:"color" binding:"omitempty,max=30"`
	Status             string     `json:"status" binding:"omitempty,oneof=active inactive in_maintenance"`
	WeightCapacityKg   *float64   `json:"weight_capacity_kg"`
	VolumeCapacityM3   *float64   `json:"volume_capacity_m3"`
	PhotoURL           string     `json:"photo_url" binding:"omitempty,max=512"`
	InsuranceExpiry    *time.Time `json:"insurance_expiry"`
	RegistrationExpiry *time.Time `json:"registration_expiry"`
}

type UpdateVehicleRequest struct {
	Type               string     `json:"type" binding:"required,oneof=motorcycle car van truck"`
	LicensePlate       string     `json:"license_plate" binding:"required,max=20"`
	Brand              string     `json:"brand" binding:"omitempty,max=64"`
	VehicleModel       string     `json:"vehicle_model" binding:"omitempty,max=64"`
	Year               *int       `json:"year"`
	Color              string     `json:"color" binding:"omitempty,max=30"`
	Status             string     `json:"status" binding:"required,oneof=active inactive in_maintenance"`
	WeightCapacityKg   *float64   `json:"weight_capacity_kg"`
	VolumeCapacityM3   *float64   `json:"volume_capacity_m3"`
	PhotoURL           string     `json:"photo_url" binding:"omitempty,max=512"`
	InsuranceExpiry    *time.Time `json:"insurance_expiry"`
	RegistrationExpiry *time.Time `json:"registration_expiry"`
}
