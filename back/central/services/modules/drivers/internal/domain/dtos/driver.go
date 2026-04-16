package dtos

import "time"

type ListDriversParams struct {
	BusinessID uint
	Search     string
	Status     string
	Page       int
	PageSize   int
}

func (p ListDriversParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type CreateDriverDTO struct {
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
}

type UpdateDriverDTO struct {
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
}
