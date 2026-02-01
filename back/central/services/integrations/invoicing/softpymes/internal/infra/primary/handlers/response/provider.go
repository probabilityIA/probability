package response

import "time"

// Provider es la respuesta de un proveedor de facturación Softpymes
type Provider struct {
	ID               uint                   `json:"id"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description,omitempty"`
	ProviderTypeCode string                 `json:"provider_type_code"`
	BusinessID       uint                   `json:"business_id"`
	Config           map[string]interface{} `json:"config,omitempty"`
	IsActive         bool                   `json:"is_active"`
	IsDefault        bool                   `json:"is_default"`
}

// ProviderList es la respuesta de listado de proveedores
type ProviderList struct {
	Items      []Provider `json:"items"`
	TotalCount int64      `json:"total_count"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
}

// ProviderType es la respuesta de un tipo de proveedor
type ProviderType struct {
	ID                 uint   `json:"id"`
	Code               string `json:"code"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Icon               string `json:"icon,omitempty"`
	ImageURL           string `json:"image_url,omitempty"`
	ApiBaseURL         string `json:"api_base_url"`
	DocumentationURL   string `json:"documentation_url,omitempty"`
	SupportedCountries string `json:"supported_countries,omitempty"`
}

// TestProviderResult es el resultado de probar la conexión con un proveedor
type TestProviderResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
