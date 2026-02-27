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
	Locations     []WarehouseLocation
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
