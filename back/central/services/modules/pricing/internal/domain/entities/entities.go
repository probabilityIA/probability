package entities

import "time"

// ClientPricingRule representa un ajuste de precio personalizado por cliente
type ClientPricingRule struct {
	ID              uint
	BusinessID      uint
	ClientID        uint
	ClientName      string  // populated by JOIN
	ProductID       *string // nil = regla global para todos los productos
	ProductName     string  // populated by JOIN
	AdjustmentType  string  // "percentage" o "fixed"
	AdjustmentValue float64 // positivo=incremento, negativo=descuento
	IsActive        bool
	Priority        int
	Description     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// QuantityDiscount representa un descuento por volumen de compra
type QuantityDiscount struct {
	ID              uint
	BusinessID      uint
	ProductID       *string // nil = aplica a todos los productos
	ProductName     string  // populated by JOIN
	MinQuantity     int
	DiscountPercent float64
	IsActive        bool
	Description     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PriceResult es el resultado del cálculo de precio
type PriceResult struct {
	BasePrice               float64
	AdjustedPrice           float64 // después de regla de cliente
	FinalPrice              float64 // después de descuento por cantidad
	ClientAdjustmentApplied bool
	ClientAdjustmentType    string
	ClientAdjustmentValue   float64
	QuantityDiscountApplied bool
	QuantityDiscountPercent float64
}
