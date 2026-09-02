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
	OverageAccepted      bool
	OverageAcceptedAt    *time.Time
	OverageAmountDue     *float64
	OverageAmountPaidAt  *time.Time
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

// ExpiringBusiness es un negocio detectado por el chequeo de vencimiento,
// junto con el codigo del plan que tiene vigente en ese momento. El codigo
// determina que mensaje de aviso le corresponde (planes pagos vs free/trial).
type ExpiringBusiness struct {
	BusinessID    uint
	PlanCode      string
	EndDate       time.Time
	CutoffDay     *int
	CourtesyUntil *time.Time
}

type BusinessSubscriptionMeta struct {
	Status        string
	EndDate       *time.Time
	CutoffDay     *int
	CourtesyUntil *time.Time
}
