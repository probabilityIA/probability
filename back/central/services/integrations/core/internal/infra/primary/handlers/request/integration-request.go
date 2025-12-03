package request

// CreateIntegrationRequest representa la solicitud para crear una integración
type CreateIntegrationRequest struct {
	Name        string                 `json:"name" binding:"required" example:"WhatsApp Principal"`
	Code        string                 `json:"code" binding:"required" example:"whatsapp_platform"`
	Type        string                 `json:"type" binding:"required" example:"whatsapp"`
	Category    string                 `json:"category" binding:"required" example:"internal"`
	BusinessID  *uint                  `json:"business_id" example:"16"` // NULL para integraciones globales
	IsActive    bool                   `json:"is_active" example:"true"`
	IsDefault   bool                   `json:"is_default" example:"true"`
	Config      map[string]interface{} `json:"config"`      // Configuración flexible
	Credentials map[string]interface{} `json:"credentials"` // Credenciales (se encriptarán)
	Description string                 `json:"description" example:"Integración principal de WhatsApp"`
}

// UpdateIntegrationRequest representa la solicitud para actualizar una integración
type UpdateIntegrationRequest struct {
	Name        *string                 `json:"name" example:"WhatsApp Actualizado"`
	Code        *string                 `json:"code" example:"whatsapp_platform"`
	IsActive    *bool                   `json:"is_active" example:"true"`
	IsDefault   *bool                   `json:"is_default" example:"true"`
	Config      *map[string]interface{} `json:"config"`      // Configuración flexible
	Credentials *map[string]interface{} `json:"credentials"` // Credenciales (se encriptarán)
	Description *string                 `json:"description" example:"Nueva descripción"`
}

// GetIntegrationsRequest representa los parámetros de consulta para obtener integraciones
type GetIntegrationsRequest struct {
	Page       int     `form:"page" example:"1"`
	PageSize   int     `form:"page_size" example:"10"`
	Type       *string `form:"type" example:"whatsapp"`
	Category   *string `form:"category" example:"internal"`
	BusinessID *uint   `form:"business_id" example:"16"`
	IsActive   *bool   `form:"is_active" example:"true"`
	Search     *string `form:"search" example:"whatsapp"`
}
