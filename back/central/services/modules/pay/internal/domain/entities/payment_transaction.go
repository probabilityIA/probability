package entities

import "time"

// PaymentStatus representa el estado de una transacción
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
)

// PaymentTransaction es la entidad de dominio para transacciones de pago
type PaymentTransaction struct {
	ID            uint
	BusinessID    uint
	Amount        float64
	Currency      string
	Status        PaymentStatus
	GatewayCode   string
	ExternalID    *string
	Reference     string
	PaymentMethod string
	Description   string
	CallbackURL   *string
	Metadata      map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
