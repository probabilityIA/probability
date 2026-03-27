package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

// ClientPricingRuleResponse respuesta de regla de precio
type ClientPricingRuleResponse struct {
	ID              uint      `json:"id"`
	BusinessID      uint      `json:"business_id"`
	ClientID        uint      `json:"client_id"`
	ClientName      string    `json:"client_name"`
	ProductID       *string   `json:"product_id"`
	ProductName     string    `json:"product_name"`
	AdjustmentType  string    `json:"adjustment_type"`
	AdjustmentValue float64   `json:"adjustment_value"`
	IsActive        bool      `json:"is_active"`
	Priority        int       `json:"priority"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// FromRuleEntity convierte una entidad de dominio a response
func FromRuleEntity(r *entities.ClientPricingRule) ClientPricingRuleResponse {
	return ClientPricingRuleResponse{
		ID:              r.ID,
		BusinessID:      r.BusinessID,
		ClientID:        r.ClientID,
		ClientName:      r.ClientName,
		ProductID:       r.ProductID,
		ProductName:     r.ProductName,
		AdjustmentType:  r.AdjustmentType,
		AdjustmentValue: r.AdjustmentValue,
		IsActive:        r.IsActive,
		Priority:        r.Priority,
		Description:     r.Description,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

// PricingRulesListResponse respuesta paginada de reglas
type PricingRulesListResponse struct {
	Data       []ClientPricingRuleResponse `json:"data"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"page_size"`
	TotalPages int                         `json:"total_pages"`
}
