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
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
