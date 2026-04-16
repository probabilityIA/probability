package request

// CreateClientPricingRuleRequest payload de creación de regla de precio
type CreateClientPricingRuleRequest struct {
	ClientID        uint    `json:"client_id" binding:"required"`
	ProductID       *string `json:"product_id" binding:"omitempty"`
	AdjustmentType  string  `json:"adjustment_type" binding:"required,oneof=percentage fixed"`
	AdjustmentValue float64 `json:"adjustment_value" binding:"required"`
	IsActive        *bool   `json:"is_active"`
	Priority        int     `json:"priority"`
	Description     string  `json:"description" binding:"omitempty,max=255"`
}

// UpdateClientPricingRuleRequest payload de actualización de regla de precio
type UpdateClientPricingRuleRequest struct {
	ClientID        uint    `json:"client_id" binding:"required"`
	ProductID       *string `json:"product_id" binding:"omitempty"`
	AdjustmentType  string  `json:"adjustment_type" binding:"required,oneof=percentage fixed"`
	AdjustmentValue float64 `json:"adjustment_value" binding:"required"`
	IsActive        *bool   `json:"is_active"`
	Priority        int     `json:"priority"`
	Description     string  `json:"description" binding:"omitempty,max=255"`
}
