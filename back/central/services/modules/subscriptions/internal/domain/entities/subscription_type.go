package entities

import "time"

type SubscriptionType struct {
	ID                   uint
	Name                 string
	Code                 string
	Description          string
	Price                float64
	BillingPeriod        string
	Active               bool
	ModuleCodes          []string
	MaxEcommerceChannels int
	BusinessID           *uint
	IncludedShipments    *int
	ShipmentOveragePrice *float64
	IncludedInvoices     *int
	InvoiceOveragePrice  *float64
	IncludedOrders       *int
	OrderOveragePrice    *float64
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
