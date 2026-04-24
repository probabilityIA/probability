package request

import "time"

type CreateLotDTO struct {
	BusinessID      uint
	ProductID       string
	LotCode         string
	ManufactureDate *time.Time
	ExpirationDate  *time.Time
	ReceivedAt      *time.Time
	SupplierID      *uint
	Status          string
}

type UpdateLotDTO struct {
	ID              uint
	BusinessID      uint
	LotCode         string
	ManufactureDate *time.Time
	ExpirationDate  *time.Time
	ReceivedAt      *time.Time
	SupplierID      *uint
	Status          string
}

type CreateSerialDTO struct {
	BusinessID   uint
	ProductID    string
	SerialNumber string
	LotID        *uint
	LocationID   *uint
	StateCode    string
}

type UpdateSerialDTO struct {
	ID         uint
	BusinessID uint
	LotID      *uint
	LocationID *uint
	StateCode  string
}

type ChangeInventoryStateDTO struct {
	BusinessID    uint
	LevelID       uint
	FromStateCode string
	ToStateCode   string
	Quantity      int
	Reason        string
	CreatedByID   *uint
}

type ConvertUoMDTO struct {
	BusinessID  uint
	ProductID   string
	FromUomCode string
	ToUomCode   string
	Quantity    float64
}

type CreateProductUoMDTO struct {
	BusinessID       uint
	ProductID        string
	UomCode          string
	ConversionFactor float64
	IsBase           bool
	Barcode          string
}
