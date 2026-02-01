package entities

import "time"

// ProbabilityPayment representa un pago que se guarda en la base de datos
// ✅ ENTIDAD PURA - SIN TAGS
type ProbabilityPayment struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	OrderID         string
	PaymentMethodID uint

	Amount       float64
	Currency     string
	ExchangeRate *float64

	Status      string
	PaidAt      *time.Time
	ProcessedAt *time.Time

	TransactionID    *string
	PaymentReference *string
	Gateway          *string

	RefundAmount  *float64
	RefundedAt    *time.Time
	FailureReason *string
	Metadata      []byte
}
