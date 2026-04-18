package entities

import "time"

type InventoryLot struct {
	ID              uint
	BusinessID      uint
	ProductID       string
	LotCode         string
	ManufactureDate *time.Time
	ExpirationDate  *time.Time
	ReceivedAt      *time.Time
	SupplierID      *uint
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time

	ProductName string
	ProductSKU  string
}
