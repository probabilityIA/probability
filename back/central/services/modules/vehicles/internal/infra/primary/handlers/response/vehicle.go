package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/entities"
)

type VehicleResponse struct {
	ID                 uint       `json:"id"`
	BusinessID         uint       `json:"business_id"`
	Type               string     `json:"type"`
	LicensePlate       string     `json:"license_plate"`
	Brand              string     `json:"brand"`
	VehicleModel       string     `json:"vehicle_model"`
	Year               *int       `json:"year"`
	Color              string     `json:"color"`
	Status             string     `json:"status"`
	WeightCapacityKg   *float64   `json:"weight_capacity_kg"`
	VolumeCapacityM3   *float64   `json:"volume_capacity_m3"`
	PhotoURL           string     `json:"photo_url"`
	InsuranceExpiry    *time.Time `json:"insurance_expiry"`
	RegistrationExpiry *time.Time `json:"registration_expiry"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type VehiclesListResponse struct {
	Data       []VehicleResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

func FromEntity(v *entities.Vehicle) VehicleResponse {
	return VehicleResponse{
		ID:                 v.ID,
		BusinessID:         v.BusinessID,
		Type:               v.Type,
		LicensePlate:       v.LicensePlate,
		Brand:              v.Brand,
		VehicleModel:       v.VehicleModel,
		Year:               v.Year,
		Color:              v.Color,
		Status:             v.Status,
		WeightCapacityKg:   v.WeightCapacityKg,
		VolumeCapacityM3:   v.VolumeCapacityM3,
		PhotoURL:           v.PhotoURL,
		InsuranceExpiry:    v.InsuranceExpiry,
		RegistrationExpiry: v.RegistrationExpiry,
		CreatedAt:          v.CreatedAt,
		UpdatedAt:          v.UpdatedAt,
	}
}
