package entities

import "time"

type Driver struct {
	ID             uint
	BusinessID     uint
	FirstName      string
	LastName       string
	Email          string
	Phone          string
	Identification string
	Status         string
	PhotoURL       string
	LicenseType    string
	LicenseExpiry  *time.Time
	WarehouseID    *uint
	Notes          *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
