package request

// CreateProvider es el request para crear un proveedor de facturación Softpymes
type CreateProvider struct {
	Name             string                 `json:"name" binding:"required,min=3,max=100"`
	ProviderTypeCode string                 `json:"provider_type_code" binding:"required"` // softpymes, siigo, etc.
	BusinessID       uint                   `json:"business_id" binding:"required"`
	Description      *string                `json:"description,omitempty"`
	Config           map[string]interface{} `json:"config"`                         // Configuración específica del proveedor
	Credentials      map[string]interface{} `json:"credentials" binding:"required"` // API keys, secrets
	IsDefault        bool                   `json:"is_default"`                     // Si es el proveedor por defecto
	CreatedByUserID  uint                   `json:"created_by_user_id" binding:"required"`
}
