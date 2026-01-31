package domain

// Constantes para estados de orden
const (
	OrderStatusAny       = "any"
	OrderStatusOpen      = "open"
	OrderStatusClosed    = "closed"
	OrderStatusCancelled = "cancelled"
)

// Constantes para estados financieros
const (
	FinancialStatusAny               = "any"
	FinancialStatusAuthorized        = "authorized"
	FinancialStatusPending           = "pending"
	FinancialStatusPaid              = "paid"
	FinancialStatusPartiallyPaid     = "partially_paid"
	FinancialStatusRefunded          = "refunded"
	FinancialStatusVoided            = "voided"
	FinancialStatusPartiallyRefunded = "partially_refunded"
	FinancialStatusUnpaid            = "unpaid"
)

// Constantes para estados de fulfillment
const (
	FulfillmentStatusAny         = "any"
	FulfillmentStatusShipped     = "shipped"
	FulfillmentStatusPartial     = "partial"
	FulfillmentStatusUnshipped   = "unshipped"
	FulfillmentStatusUnfulfilled = "unfulfilled"
)
