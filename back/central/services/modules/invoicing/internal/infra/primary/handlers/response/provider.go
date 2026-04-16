package response

import "time"

// Provider es la respuesta de un proveedor de facturación
type Provider struct {
	ID               uint                   `json:"id"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Name             string                 `json:"name"`
	ProviderTypeCode string                 `json:"provider_type_code"`
	BusinessID       uint                   `json:"business_id"`
	Config           map[string]interface{} `json:"config,omitempty"`
	Credentials      map[string]interface{} `json:"credentials,omitempty"` // Credenciales ofuscadas (no enviar keys completos)
	IsActive         bool                   `json:"is_active"`
}

// ProviderList es la respuesta de listado de proveedores
type ProviderList struct {
	Items      []Provider `json:"items"`
	TotalCount int64      `json:"total_count"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
}

// TestProviderResult es el resultado de probar la conexión con un proveedor
type TestProviderResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
