package dtos

type CreateSubscriptionTypeDTO struct {
	Name          string
	Code          string
	Description   string
	Price         float64
	BillingPeriod string
	ModuleCodes   []string
}

type UpdateSubscriptionTypeDTO struct {
	ID            uint
	Name          string
	Description   string
	Price         float64
	BillingPeriod string
	Active        bool
	ModuleCodes   []string
}
