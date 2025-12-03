package domain

import "time"

// ───────────────────────────────────────────
//
//	PAYMENT METHOD DTOs
//
// ───────────────────────────────────────────

// CreatePaymentMethodRequest representa la solicitud para crear un método de pago
type CreatePaymentMethodRequest struct {
	Code        string `json:"code" binding:"required,max=64"`
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"required,oneof=card digital_wallet bank_transfer cash"`
	Provider    string `json:"provider" binding:"max=64"`
	Icon        string `json:"icon" binding:"max=255"`
	Color       string `json:"color" binding:"max=32"`
}

// UpdatePaymentMethodRequest representa la solicitud para actualizar un método de pago
type UpdatePaymentMethodRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"required,oneof=card digital_wallet bank_transfer cash"`
	Provider    string `json:"provider" binding:"max=64"`
	Icon        string `json:"icon" binding:"max=255"`
	Color       string `json:"color" binding:"max=32"`
}

// PaymentMethodResponse representa la respuesta de un método de pago
type PaymentMethodResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Provider    string    `json:"provider"`
	IsActive    bool      `json:"is_active"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PaymentMethodsListResponse representa la respuesta paginada de métodos de pago
type PaymentMethodsListResponse struct {
	Data       []PaymentMethodResponse `json:"data"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalPages int                     `json:"total_pages"`
}

// ───────────────────────────────────────────
//
//	PAYMENT METHOD MAPPING DTOs
//
// ───────────────────────────────────────────

// CreatePaymentMappingRequest representa la solicitud para crear un mapeo
type CreatePaymentMappingRequest struct {
	IntegrationType string `json:"integration_type" binding:"required,oneof=shopify whatsapp mercadolibre"`
	OriginalMethod  string `json:"original_method" binding:"required,max=128"`
	PaymentMethodID uint   `json:"payment_method_id" binding:"required"`
	Priority        int    `json:"priority"`
}

// UpdatePaymentMappingRequest representa la solicitud para actualizar un mapeo
type UpdatePaymentMappingRequest struct {
	OriginalMethod  string `json:"original_method" binding:"required,max=128"`
	PaymentMethodID uint   `json:"payment_method_id" binding:"required"`
	Priority        int    `json:"priority"`
}

// PaymentMappingResponse representa la respuesta de un mapeo
type PaymentMappingResponse struct {
	ID              uint                  `json:"id"`
	IntegrationType string                `json:"integration_type"`
	OriginalMethod  string                `json:"original_method"`
	PaymentMethodID uint                  `json:"payment_method_id"`
	PaymentMethod   PaymentMethodResponse `json:"payment_method"`
	IsActive        bool                  `json:"is_active"`
	Priority        int                   `json:"priority"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

// PaymentMappingsListResponse representa la respuesta de lista de mapeos
type PaymentMappingsListResponse struct {
	Data  []PaymentMappingResponse `json:"data"`
	Total int64                    `json:"total"`
}

// PaymentMappingsByIntegrationResponse agrupa mapeos por tipo de integración
type PaymentMappingsByIntegrationResponse struct {
	IntegrationType string                   `json:"integration_type"`
	Mappings        []PaymentMappingResponse `json:"mappings"`
}
