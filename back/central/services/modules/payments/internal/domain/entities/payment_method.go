package entities

import "time"

// PaymentMethod representa un método de pago en el dominio (PURO - sin tags)
type PaymentMethod struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string // card, digital_wallet, bank_transfer, cash
	Provider    string
	IsActive    bool
	Icon        string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
