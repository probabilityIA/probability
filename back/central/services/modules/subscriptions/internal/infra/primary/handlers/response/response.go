package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

type SubscriptionTypeResponse struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Code                 string    `json:"code"`
	Description          string    `json:"description"`
	Price                float64   `json:"price"`
	BillingPeriod        string    `json:"billing_period"`
	Active               bool      `json:"active"`
	ModuleCodes          []string  `json:"module_codes"`
	MaxEcommerceChannels int       `json:"max_ecommerce_channels"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func FromSubscriptionType(t *entities.SubscriptionType) SubscriptionTypeResponse {
	return SubscriptionTypeResponse{
		ID:                   t.ID,
		Name:                 t.Name,
		Code:                 t.Code,
		Description:          t.Description,
		Price:                t.Price,
		BillingPeriod:        t.BillingPeriod,
		Active:               t.Active,
		ModuleCodes:          t.ModuleCodes,
		MaxEcommerceChannels: t.MaxEcommerceChannels,
		CreatedAt:            t.CreatedAt,
		UpdatedAt:            t.UpdatedAt,
	}
}

func FromSubscriptionTypes(types []entities.SubscriptionType) []SubscriptionTypeResponse {
	result := make([]SubscriptionTypeResponse, len(types))
	for i, t := range types {
		result[i] = FromSubscriptionType(&t)
	}
	return result
}

type SubscriptionResponse struct {
	ID                   uint      `json:"id"`
	BusinessID           uint      `json:"business_id"`
	SubscriptionTypeID   uint      `json:"subscription_type_id"`
	SubscriptionTypeName string    `json:"subscription_type_name"`
	Months               int       `json:"months"`
	Amount               float64   `json:"amount"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Status               string    `json:"status"`
	PaymentReference     *string   `json:"payment_reference,omitempty"`
	Notes                *string   `json:"notes,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

func FromSubscription(s *entities.BusinessSubscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:                   s.ID,
		BusinessID:           s.BusinessID,
		SubscriptionTypeID:   s.SubscriptionTypeID,
		SubscriptionTypeName: s.SubscriptionTypeName,
		Months:               s.Months,
		Amount:               s.Amount,
		StartDate:            s.StartDate,
		EndDate:              s.EndDate,
		Status:               s.Status,
		PaymentReference:     s.PaymentReference,
		Notes:                s.Notes,
		CreatedAt:            s.CreatedAt,
	}
}

type OverrideResponse struct {
	ID              uint      `json:"id"`
	BusinessID      uint      `json:"business_id"`
	ModuleCode      string    `json:"module_code"`
	GrantedByUserID uint      `json:"granted_by_user_id"`
	Notes           *string   `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

func FromOverride(o *entities.BusinessModuleOverride) OverrideResponse {
	return OverrideResponse{
		ID:              o.ID,
		BusinessID:      o.BusinessID,
		ModuleCode:      o.ModuleCode,
		GrantedByUserID: o.GrantedByUserID,
		Notes:           o.Notes,
		CreatedAt:       o.CreatedAt,
	}
}

func FromOverrides(overrides []entities.BusinessModuleOverride) []OverrideResponse {
	result := make([]OverrideResponse, len(overrides))
	for i, o := range overrides {
		result[i] = FromOverride(&o)
	}
	return result
}
