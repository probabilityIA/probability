package request

import "time"

type CreateLotBody struct {
	ProductID       string     `json:"product_id" binding:"required"`
	LotCode         string     `json:"lot_code" binding:"required,min=1,max=100"`
	ManufactureDate *time.Time `json:"manufacture_date"`
	ExpirationDate  *time.Time `json:"expiration_date"`
	ReceivedAt      *time.Time `json:"received_at"`
	SupplierID      *uint      `json:"supplier_id"`
	Status          string     `json:"status"`
}

type UpdateLotBody struct {
	LotCode         string     `json:"lot_code"`
	ManufactureDate *time.Time `json:"manufacture_date"`
	ExpirationDate  *time.Time `json:"expiration_date"`
	ReceivedAt      *time.Time `json:"received_at"`
	SupplierID      *uint      `json:"supplier_id"`
	Status          string     `json:"status"`
}

type CreateSerialBody struct {
	ProductID    string `json:"product_id" binding:"required"`
	SerialNumber string `json:"serial_number" binding:"required,min=1,max=100"`
	LotID        *uint  `json:"lot_id"`
	LocationID   *uint  `json:"location_id"`
	StateCode    string `json:"state_code"`
}

type UpdateSerialBody struct {
	LotID      *uint  `json:"lot_id"`
	LocationID *uint  `json:"location_id"`
	StateCode  string `json:"state_code"`
}

type ChangeInventoryStateBody struct {
	LevelID       uint   `json:"level_id" binding:"required,min=1"`
	FromStateCode string `json:"from_state_code" binding:"required"`
	ToStateCode   string `json:"to_state_code" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required,min=1"`
	Reason        string `json:"reason"`
}

type CreateProductUoMBody struct {
	UomCode          string  `json:"uom_code" binding:"required"`
	ConversionFactor float64 `json:"conversion_factor"`
	IsBase           bool    `json:"is_base"`
	Barcode          string  `json:"barcode"`
}

type ConvertUoMBody struct {
	ProductID   string  `json:"product_id" binding:"required"`
	FromUomCode string  `json:"from_uom_code" binding:"required"`
	ToUomCode   string  `json:"to_uom_code" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required"`
}
