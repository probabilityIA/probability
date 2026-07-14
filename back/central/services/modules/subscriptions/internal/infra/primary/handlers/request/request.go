package request

type CreateSubscriptionTypeRequest struct {
	Name          string   `json:"name" binding:"required"`
	Code          string   `json:"code" binding:"required"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required,gt=0"`
	BillingPeriod string   `json:"billing_period"`
	ModuleCodes   []string `json:"module_codes"`
}

type UpdateSubscriptionTypeRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required,gt=0"`
	BillingPeriod string   `json:"billing_period"`
	Active        bool     `json:"active"`
	ModuleCodes   []string `json:"module_codes"`
}

type PurchaseSubscriptionRequest struct {
	SubscriptionTypeID uint `json:"subscription_type_id" binding:"required"`
	Months             int  `json:"months" binding:"required,gt=0"`
}

type RegisterPaymentRequest struct {
	BusinessID         uint    `json:"business_id" binding:"required"`
	SubscriptionTypeID uint    `json:"subscription_type_id" binding:"required"`
	Months             int     `json:"months" binding:"required,gt=0"`
	PaymentReference   *string `json:"payment_reference"`
	Notes              *string `json:"notes"`
}

type GrantOverrideRequest struct {
	BusinessID uint    `json:"business_id" binding:"required"`
	ModuleCode string  `json:"module_code" binding:"required"`
	Notes      *string `json:"notes"`
}
