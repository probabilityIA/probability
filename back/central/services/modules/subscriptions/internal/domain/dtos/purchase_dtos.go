package dtos

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
	PaymentReference   *string
	Notes              *string
}
