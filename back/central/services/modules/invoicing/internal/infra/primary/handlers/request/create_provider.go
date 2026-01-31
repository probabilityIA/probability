package request

// CreateProvider es el request para crear un proveedor de facturación
type CreateProvider struct {
	Name             string                 `json:"name" binding:"required,min=3,max=100"`
	ProviderTypeCode string                 `json:"provider_type_code" binding:"required"` // softpymes, siigo, etc.
	BusinessID       uint                   `json:"business_id" binding:"required"`
	Config           map[string]interface{} `json:"config"`                                // Configuración específica del proveedor
	Credentials      map[string]interface{} `json:"credentials" binding:"required"`        // API keys, secrets (se encriptarán)
	IsActive         *bool                  `json:"is_active"`                             // Por defecto true
}
