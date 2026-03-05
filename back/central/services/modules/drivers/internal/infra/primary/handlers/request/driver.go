package request

import "time"

type CreateDriverRequest struct {
	FirstName      string     `json:"first_name" binding:"required,min=2,max=128"`
	LastName       string     `json:"last_name" binding:"required,min=2,max=128"`
	Email          string     `json:"email" binding:"omitempty,email,max=255"`
	Phone          string     `json:"phone" binding:"required,max=50"`
	Identification string     `json:"identification" binding:"required,max=50"`
	Status         string     `json:"status" binding:"omitempty,oneof=active inactive"`
	PhotoURL       string     `json:"photo_url" binding:"omitempty,max=512"`
	LicenseType    string     `json:"license_type" binding:"omitempty,max=20"`
	LicenseExpiry  *time.Time `json:"license_expiry"`
	WarehouseID    *uint      `json:"warehouse_id"`
	Notes          *string    `json:"notes"`
}

type UpdateDriverRequest struct {
	FirstName      string     `json:"first_name" binding:"required,min=2,max=128"`
	LastName       string     `json:"last_name" binding:"required,min=2,max=128"`
	Email          string     `json:"email" binding:"omitempty,email,max=255"`
	Phone          string     `json:"phone" binding:"required,max=50"`
	Identification string     `json:"identification" binding:"required,max=50"`
	Status         string     `json:"status" binding:"required,oneof=active inactive on_route"`
	PhotoURL       string     `json:"photo_url" binding:"omitempty,max=512"`
	LicenseType    string     `json:"license_type" binding:"omitempty,max=20"`
	LicenseExpiry  *time.Time `json:"license_expiry"`
	WarehouseID    *uint      `json:"warehouse_id"`
	Notes          *string    `json:"notes"`
}
