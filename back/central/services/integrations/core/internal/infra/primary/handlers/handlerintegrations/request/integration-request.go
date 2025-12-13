package request

// CreateIntegrationRequest representa la solicitud para crear una integración
type CreateIntegrationRequest struct {
	Name              string                 `json:"name" binding:"required" example:"WhatsApp Principal"`
	Code              string                 `json:"code" binding:"required" example:"whatsapp_platform"`
	IntegrationTypeID uint                   `json:"integration_type_id" binding:"required" example:"1"` // ID del tipo de integración
	Category          string                 `json:"category" binding:"omitempty" example:"internal"`
	BusinessID        *uint                  `json:"business_id" example:"16"` // NULL para integraciones globales
	IsActive          bool                   `json:"is_active" example:"true"`
	IsDefault         bool                   `json:"is_default" example:"true"`
	Config            map[string]interface{} `json:"config"`      // Configuración flexible
	Credentials       map[string]interface{} `json:"credentials"` // Credenciales (se encriptarán)
	Description       string                 `json:"description" example:"Integración principal de WhatsApp"`
}

// UpdateIntegrationRequest representa la solicitud para actualizar una integración
type UpdateIntegrationRequest struct {
	Name              *string                 `json:"name" example:"WhatsApp Actualizado"`
	Code              *string                 `json:"code" example:"whatsapp_platform"`
	IntegrationTypeID *uint                   `json:"integration_type_id" example:"1"` // ID del tipo de integración (opcional)
	IsActive          *bool                   `json:"is_active" example:"true"`
	IsDefault         *bool                   `json:"is_default" example:"true"`
	Config            *map[string]interface{} `json:"config"`      // Configuración flexible
	Credentials       *map[string]interface{} `json:"credentials"` // Credenciales (se encriptarán)
	Description       *string                 `json:"description" example:"Nueva descripción"`
}

// GetIntegrationsRequest representa los parámetros de consulta para obtener integraciones
type GetIntegrationsRequest struct {
	Page                int     `form:"page" example:"1"`
	PageSize            int     `form:"page_size" example:"10"`
	IntegrationTypeID   *uint   `form:"integration_type_id" example:"1"`          // Filtrar por ID del tipo de integración
	IntegrationTypeCode *string `form:"integration_type_code" example:"whatsapp"` // Filtrar por código del tipo de integración
	Category            *string `form:"category" example:"internal"`
	BusinessID          *uint   `form:"business_id" example:"16"`
	IsActive            *bool   `form:"is_active" example:"true"`
	Search              *string `form:"search" example:"whatsapp"`
}
