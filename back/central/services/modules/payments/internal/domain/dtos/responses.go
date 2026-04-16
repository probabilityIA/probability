package dtos

import "time"

// PaymentMethodResponse representa la respuesta de un método de pago (PURO - sin tags)
type PaymentMethodResponse struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Provider    string
	IsActive    bool
	Icon        string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PaymentMethodsListResponse representa la respuesta paginada de métodos de pago
type PaymentMethodsListResponse struct {
	Data       []PaymentMethodResponse
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// PaymentMappingResponse representa la respuesta de un mapeo
type PaymentMappingResponse struct {
	ID              uint
	IntegrationType string
	OriginalMethod  string
	PaymentMethodID uint
	PaymentMethod   PaymentMethodResponse
	IsActive        bool
	Priority        int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PaymentMappingsListResponse representa la respuesta de lista de mapeos
type PaymentMappingsListResponse struct {
	Data  []PaymentMappingResponse
	Total int64
}

// PaymentMappingsByIntegrationResponse agrupa mapeos por tipo de integración
type PaymentMappingsByIntegrationResponse struct {
	IntegrationType string
	Mappings        []PaymentMappingResponse
}
