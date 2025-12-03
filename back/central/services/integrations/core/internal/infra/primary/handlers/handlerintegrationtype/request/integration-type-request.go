package request

// CreateIntegrationTypeRequest representa la solicitud para crear un tipo de integración
type CreateIntegrationTypeRequest struct {
	Name              string                 `json:"name" binding:"required" example:"WhatsApp"`
	Code              string                 `json:"code" example:"whatsapp"` // Opcional, se genera automáticamente
	Description       string                 `json:"description" example:"Integración con WhatsApp Cloud API"`
	Icon              string                 `json:"icon" example:"whatsapp-icon"`
	Category          string                 `json:"category" binding:"required" example:"internal"`
	IsActive          bool                   `json:"is_active" example:"true"`
	ConfigSchema      map[string]interface{} `json:"credentials_schema"` // JSON schema para credenciales
}

// UpdateIntegrationTypeRequest representa la solicitud para actualizar un tipo de integración
type UpdateIntegrationTypeRequest struct {
	Name              *string                 `json:"name" example:"WhatsApp Actualizado"`
	Code              *string                 `json:"code" example:"whatsapp"`
	Description       *string                 `json:"description" example:"Nueva descripción"`
	Icon              *string                 `json:"icon" example:"whatsapp-icon"`
	Category          *string                 `json:"category" example:"internal"`
	IsActive          *bool                   `json:"is_active" example:"true"`
	ConfigSchema      *map[string]interface{} `json:"credentials_schema"`
}
