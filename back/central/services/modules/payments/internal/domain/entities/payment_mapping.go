package entities

import "time"

// PaymentMethodMapping representa el mapeo entre métodos de pago de integraciones (PURO - sin tags)
type PaymentMethodMapping struct {
	ID              uint
	IntegrationType string // shopify, whatsapp, mercadolibre
	OriginalMethod  string
	PaymentMethodID uint
	PaymentMethod   PaymentMethod
	IsActive        bool
	Priority        int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
