package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/entities"
)

type DriverResponse struct {
	ID             uint       `json:"id"`
	BusinessID     uint       `json:"business_id"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	Identification string     `json:"identification"`
	Status         string     `json:"status"`
	PhotoURL       string     `json:"photo_url"`
	LicenseType    string     `json:"license_type"`
	LicenseExpiry  *time.Time `json:"license_expiry"`
	WarehouseID    *uint      `json:"warehouse_id"`
	Notes          *string    `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type DriversListResponse struct {
	Data       []DriverResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

func FromEntity(d *entities.Driver) DriverResponse {
	return DriverResponse{
		ID:             d.ID,
		BusinessID:     d.BusinessID,
		FirstName:      d.FirstName,
		LastName:       d.LastName,
		Email:          d.Email,
		Phone:          d.Phone,
		Identification: d.Identification,
		Status:         d.Status,
		PhotoURL:       d.PhotoURL,
		LicenseType:    d.LicenseType,
		LicenseExpiry:  d.LicenseExpiry,
		WarehouseID:    d.WarehouseID,
		Notes:          d.Notes,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}
