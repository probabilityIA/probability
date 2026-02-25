package request

import "mime/multipart"

// CreateIntegrationTypeRequest representa la solicitud para crear un tipo de integración
// Soporta tanto JSON como multipart/form-data
type CreateIntegrationTypeRequest struct {
	Name                string                 `json:"name" form:"name" binding:"required" example:"WhatsApp"`
	Code                string                 `json:"code" form:"code" example:"whatsapp"` // Opcional, se genera automáticamente
	Description         string                 `json:"description" form:"description" example:"Integración con WhatsApp Cloud API"`
	Icon                string                 `json:"icon" form:"icon" example:"whatsapp-icon"`
	CategoryID          uint                   `json:"category_id" form:"category_id" binding:"required" example:"1"`
	IsActive            bool                   `json:"is_active" form:"is_active" example:"true"`
	InDevelopment       bool                   `json:"in_development" form:"in_development" example:"false"`
	ConfigSchema        map[string]interface{} `json:"credentials_schema" form:"credentials_schema"` // JSON schema para credenciales
	ImageFile           *multipart.FileHeader  `form:"image_file"`                                   // Archivo de imagen para subir a S3
	BaseURL             string                 `json:"base_url" form:"base_url"`
	BaseURLTest         string                 `json:"base_url_test" form:"base_url_test"`
	PlatformCredentials map[string]interface{} `json:"platform_credentials" form:"platform_credentials"` // Credenciales de plataforma (se encriptarán)
}

// UpdateIntegrationTypeRequest representa la solicitud para actualizar un tipo de integración
// Soporta tanto JSON como multipart/form-data
type UpdateIntegrationTypeRequest struct {
	Name                *string                 `json:"name" form:"name" example:"WhatsApp Actualizado"`
	Code                *string                 `json:"code" form:"code" example:"whatsapp"`
	Description         *string                 `json:"description" form:"description" example:"Nueva descripción"`
	Icon                *string                 `json:"icon" form:"icon" example:"whatsapp-icon"`
	CategoryID          *uint                   `json:"category_id" form:"category_id" example:"1"`
	IsActive            *bool                   `json:"is_active" form:"is_active" example:"true"`
	InDevelopment       *bool                   `json:"in_development" form:"in_development" example:"false"`
	ConfigSchema        *map[string]interface{} `json:"credentials_schema" form:"credentials_schema"`
	ImageFile           *multipart.FileHeader   `form:"image_file"`                       // Archivo de imagen para subir a S3
	RemoveImage         *bool                   `json:"remove_image" form:"remove_image"` // Flag para eliminar la imagen existente
	BaseURL             *string                 `json:"base_url" form:"base_url"`
	BaseURLTest         *string                 `json:"base_url_test" form:"base_url_test"`
	PlatformCredentials *map[string]interface{} `json:"platform_credentials" form:"platform_credentials"` // Credenciales de plataforma (se encriptarán)
}
