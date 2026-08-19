package dtos

type CreateSubscriptionTypeDTO struct {
	Name                 string
	Code                 string
	Description          string
	Price                float64
	BillingPeriod        string
	ModuleCodes          []string
	MaxEcommerceChannels int
	BusinessID           *uint
	IncludedShipments    *int
	ShipmentOveragePrice *float64
	IncludedInvoices     *int
	InvoiceOveragePrice  *float64
	IncludedOrders       *int
	OrderOveragePrice    *float64
}

type CreateCustomPlanDTO struct {
	Name                 string
	Code                 string
	Description          string
	Price                float64
	BillingPeriod        string
	ModuleCodes          []string
	MaxEcommerceChannels int
	BusinessID           uint
	Months               int
	PaymentReference     *string
	Notes                *string
	IncludedShipments    *int
	ShipmentOveragePrice *float64
	IncludedInvoices     *int
	InvoiceOveragePrice  *float64
	IncludedOrders       *int
	OrderOveragePrice    *float64
}

type UpdateSubscriptionTypeDTO struct {
	ID                   uint
	Name                 string
	Description          string
	Price                float64
	BillingPeriod        string
	Active               bool
	ModuleCodes          []string
	MaxEcommerceChannels int
	IncludedShipments    *int
	ShipmentOveragePrice *float64
	IncludedInvoices     *int
	InvoiceOveragePrice  *float64
	IncludedOrders       *int
	OrderOveragePrice    *float64
}
