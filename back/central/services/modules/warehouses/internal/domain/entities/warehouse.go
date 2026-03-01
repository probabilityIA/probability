package entities

import "time"

// Warehouse representa una bodega en el dominio
type Warehouse struct {
	ID            uint
	BusinessID    uint
	Name          string
	Code          string
	Address       string
	City          string
	State         string
	Country       string
	ZipCode       string
	Phone         string
	ContactName   string
	ContactEmail  string
	IsActive      bool
	IsDefault     bool
	IsFulfillment bool
	Company       string
	FirstName     string
	LastName      string
	Email         string
	Suburb        string
	CityDaneCode  string
	PostalCode    string
	Street        string
	Latitude      *float64
	Longitude     *float64
	Locations     []WarehouseLocation
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
