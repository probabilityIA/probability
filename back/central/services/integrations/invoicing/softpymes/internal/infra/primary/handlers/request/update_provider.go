package request

// UpdateProvider es el request para actualizar un proveedor de facturación
type UpdateProvider struct {
	Name        *string                 `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	Description *string                 `json:"description,omitempty"`
	Config      *map[string]interface{} `json:"config,omitempty"`
	Credentials *map[string]interface{} `json:"credentials,omitempty"`
	IsActive    *bool                   `json:"is_active,omitempty"`
	IsDefault   *bool                   `json:"is_default,omitempty"`
}
