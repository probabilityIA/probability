package dtos

import "time"

type PurchaseSubscriptionDTO struct {
	BusinessID         uint
	SubscriptionTypeID uint
	Months             int
	UserID             uint
}

type RegisterPaymentDTO struct {
	BusinessID         uint
	SubscriptionTypeID uint
	Months             int
	PaymentMethod      *string
	PaymentReference   *string
	Notes              *string
	StartDate          *time.Time
}

type EditSubscriptionDatesDTO struct {
	BusinessID uint
	StartDate  time.Time
	EndDate    time.Time
}

type ExtendCourtesyDTO struct {
	BusinessID uint
	Days       int
	Reason     string
}
