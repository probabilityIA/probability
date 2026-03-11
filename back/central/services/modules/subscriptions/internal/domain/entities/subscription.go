package entities

import "time"

// BusinessSubscription representa una suscripción o pago de un cliente en el sistema
type BusinessSubscription struct {
	ID               uint
	BusinessID       uint
	Amount           float64
	StartDate        time.Time
	EndDate          time.Time
	Status           string // 'paid', 'pending', 'rejected'
	PaymentReference *string
	Notes            *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// BusinessSubscriptionStatus representa los posibles estados de una suscripción
const (
	SubscriptionStatusPaid     = "paid"
	SubscriptionStatusPending  = "pending"
	SubscriptionStatusRejected = "rejected"
)

// BusinessStatus representa los posibles estados de bloqueo del negocio por pago
const (
	BusinessStatusActive    = "active"
	BusinessStatusExpired   = "expired"
	BusinessStatusCancelled = "cancelled"
)
