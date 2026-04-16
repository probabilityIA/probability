package response

import "time"

// PaymentMapping representa la respuesta HTTP de un mapeo
type PaymentMapping struct {
	ID              uint          `json:"id"`
	IntegrationType string        `json:"integration_type"`
	OriginalMethod  string        `json:"original_method"`
	PaymentMethodID uint          `json:"payment_method_id"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	IsActive        bool          `json:"is_active"`
	Priority        int           `json:"priority"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// PaymentMappingsList representa la respuesta HTTP de lista de mapeos
type PaymentMappingsList struct {
	Data  []PaymentMapping `json:"data"`
	Total int64            `json:"total"`
}

// PaymentMappingsByIntegration agrupa mapeos por tipo de integración (respuesta HTTP)
type PaymentMappingsByIntegration struct {
	IntegrationType string           `json:"integration_type"`
	Mappings        []PaymentMapping `json:"mappings"`
}
