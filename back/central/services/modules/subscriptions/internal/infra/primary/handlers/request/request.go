package request

type CreateSubscriptionTypeRequest struct {
	Name                 string   `json:"name" binding:"required"`
	Code                 string   `json:"code" binding:"required"`
	Description          string   `json:"description"`
	Price                float64  `json:"price" binding:"required,gt=0"`
	BillingPeriod        string   `json:"billing_period"`
	ModuleCodes          []string `json:"module_codes"`
	MaxEcommerceChannels int      `json:"max_ecommerce_channels"`
	IncludedShipments    *int     `json:"included_shipments"`
	ShipmentOveragePrice *float64 `json:"shipment_overage_price"`
	IncludedInvoices     *int     `json:"included_invoices"`
	InvoiceOveragePrice  *float64 `json:"invoice_overage_price"`
	IncludedOrders       *int     `json:"included_orders"`
	OrderOveragePrice    *float64 `json:"order_overage_price"`
}

type UpdateSubscriptionTypeRequest struct {
	Name                 string   `json:"name" binding:"required"`
	Description          string   `json:"description"`
	Price                float64  `json:"price" binding:"required,gt=0"`
	BillingPeriod        string   `json:"billing_period"`
	Active               bool     `json:"active"`
	ModuleCodes          []string `json:"module_codes"`
	MaxEcommerceChannels int      `json:"max_ecommerce_channels"`
	IncludedShipments    *int     `json:"included_shipments"`
	ShipmentOveragePrice *float64 `json:"shipment_overage_price"`
	IncludedInvoices     *int     `json:"included_invoices"`
	InvoiceOveragePrice  *float64 `json:"invoice_overage_price"`
	IncludedOrders       *int     `json:"included_orders"`
	OrderOveragePrice    *float64 `json:"order_overage_price"`
}

type CreateCustomPlanRequest struct {
	Name                 string   `json:"name" binding:"required"`
	Code                 string   `json:"code" binding:"required"`
	Description          string   `json:"description"`
	Price                float64  `json:"price" binding:"gte=0"`
	BillingPeriod        string   `json:"billing_period"`
	ModuleCodes          []string `json:"module_codes"`
	MaxEcommerceChannels int      `json:"max_ecommerce_channels"`
	BusinessID           uint     `json:"business_id" binding:"required"`
	Months               int      `json:"months" binding:"required,gt=0"`
	PaymentReference     *string  `json:"payment_reference"`
	Notes                *string  `json:"notes"`
	IncludedShipments    *int     `json:"included_shipments"`
	ShipmentOveragePrice *float64 `json:"shipment_overage_price"`
	IncludedInvoices     *int     `json:"included_invoices"`
	InvoiceOveragePrice  *float64 `json:"invoice_overage_price"`
	IncludedOrders       *int     `json:"included_orders"`
	OrderOveragePrice    *float64 `json:"order_overage_price"`
}

type PurchaseSubscriptionRequest struct {
	SubscriptionTypeID uint `json:"subscription_type_id" binding:"required"`
	Months             int  `json:"months" binding:"required,gt=0"`
}

type RegisterPaymentRequest struct {
	BusinessID         uint    `json:"business_id" binding:"required"`
	SubscriptionTypeID uint    `json:"subscription_type_id" binding:"required"`
	Months             int     `json:"months" binding:"required,gt=0"`
	PaymentMethod      *string `json:"payment_method"`
	PaymentReference   *string `json:"payment_reference"`
	Notes              *string `json:"notes"`
	StartDate          *string `json:"start_date"`
}

type EditSubscriptionDatesRequest struct {
	BusinessID uint   `json:"business_id" binding:"required"`
	StartDate  string `json:"start_date" binding:"required"`
	EndDate    string `json:"end_date" binding:"required"`
}

type GrantOverrideRequest struct {
	BusinessID uint    `json:"business_id" binding:"required"`
	ModuleCode string  `json:"module_code" binding:"required"`
	Notes      *string `json:"notes"`
	ExpiresAt  *string `json:"expires_at"`
}

type ReactivateSubscriptionRequest struct {
	BusinessID uint `json:"business_id" binding:"required"`
}

type ExtendCourtesyRequest struct {
	BusinessID uint   `json:"business_id" binding:"required"`
	Days       int    `json:"days" binding:"required,gt=0"`
	Reason     string `json:"reason" binding:"required"`
}
