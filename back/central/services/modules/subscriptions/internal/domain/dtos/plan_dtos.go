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
}
