package entities

import "time"

type BusinessSubscription struct {
	ID                   uint
	BusinessID           uint
	SubscriptionTypeID   uint
	SubscriptionTypeName string
	Months               int
	Amount               float64
	StartDate            time.Time
	EndDate              time.Time
	Status               string
	PaymentMethod        string
	PaymentReference     *string
	Notes                *string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

const (
	SubscriptionStatusPaid     = "paid"
	SubscriptionStatusPending  = "pending"
	SubscriptionStatusRejected = "rejected"
	SubscriptionStatusReverted = "reverted"
)

const (
	PaymentMethodWallet   = "WALLET"
	PaymentMethodManual   = "MANUAL"
	PaymentMethodCourtesy = "COURTESY"
)

const (
	BusinessStatusActive    = "active"
	BusinessStatusExpired   = "expired"
	BusinessStatusCancelled = "cancelled"
)
