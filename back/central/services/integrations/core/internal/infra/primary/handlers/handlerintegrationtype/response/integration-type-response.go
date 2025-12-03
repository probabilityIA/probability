package response

import "gorm.io/datatypes"

// IntegrationTypeResponse representa la respuesta de un tipo de integración
type IntegrationTypeResponse struct {
	ID                uint           `json:"id" example:"1"`
	Name              string         `json:"name" example:"WhatsApp"`
	Code              string         `json:"code" example:"whatsapp"`
	Description       string         `json:"description" example:"Integración con WhatsApp Cloud API"`
	Icon              string         `json:"icon" example:"whatsapp-icon"`
	Category          string         `json:"category" example:"internal"`
	IsActive          bool           `json:"is_active" example:"true"`
	ConfigSchema      datatypes.JSON `json:"config_schema"`
	CredentialsSchema datatypes.JSON `json:"credentials_schema"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
}

// IntegrationTypeListResponse representa la respuesta de lista de tipos de integración
type IntegrationTypeListResponse struct {
	Success bool                      `json:"success" example:"true"`
	Message string                    `json:"message" example:"Tipos de integración obtenidos exitosamente"`
	Data    []IntegrationTypeResponse `json:"data"`
}

// IntegrationTypeDetailResponse representa la respuesta de detalle de un tipo de integración
type IntegrationTypeDetailResponse struct {
	Success bool                    `json:"success" example:"true"`
	Message string                  `json:"message" example:"Tipo de integración obtenido exitosamente"`
	Data    IntegrationTypeResponse `json:"data"`
}

// IntegrationErrorResponse representa la respuesta de error (compartido con handlerintegrations)
//
//	@Description	Respuesta de error para operaciones fallidas
type IntegrationErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Error al procesar la solicitud"`
	Error   string `json:"error,omitempty" example:"Detalles del error"`
}

// IntegrationMessageResponse representa la respuesta de mensaje (compartido con handlerintegrations)
//
//	@Description	Respuesta de mensaje para operaciones simples
type IntegrationMessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operación realizada exitosamente"`
}
