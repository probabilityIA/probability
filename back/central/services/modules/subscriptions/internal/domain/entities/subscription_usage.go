package entities

import "time"

type SubscriptionUsage struct {
	PlanCode             string
	PlanName             string
	PlanPrice            float64
	BillingPeriod        string
	ModuleCodes          []string
	MaxEcommerceChannels int
	CycleStartDate       time.Time
	CycleEndDate         time.Time

	IncludedShipments    *int
	ShipmentOveragePrice *float64
	ShipmentsUsed        int64
	OverageAccepted      bool
	OverageAmountDue     *float64
	OverageAmountPaidAt  *time.Time

	IncludedInvoices    *int
	InvoiceOveragePrice *float64
	InvoicesUsed        int64

	IncludedOrders    *int
	OrderOveragePrice *float64
	OrdersUsed        int64

	ForecastedPayment *float64
}
