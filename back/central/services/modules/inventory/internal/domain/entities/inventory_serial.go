package entities

import "time"

type InventorySerial struct {
	ID                uint
	BusinessID        uint
	ProductID         string
	SerialNumber      string
	LotID             *uint
	CurrentLocationID *uint
	CurrentStateID    *uint
	ReceivedAt        *time.Time
	SoldAt            *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
