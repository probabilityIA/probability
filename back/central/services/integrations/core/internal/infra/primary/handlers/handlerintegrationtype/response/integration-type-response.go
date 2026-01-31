package response

import "gorm.io/datatypes"

// IntegrationCategoryResponse representa la respuesta de una categoría de integración
type IntegrationCategoryResponse struct {
	ID          uint   `json:"id" example:"1"`
	Code        string `json:"code" example:"ecommerce"`
	Name        string `json:"name" example:"E-commerce"`
	Description string `json:"description" example:"Plataformas de venta online"`
	Icon        string `json:"icon" example:"shopping-cart"`
	Color       string `json:"color" example:"#3B82F6"`
}

// IntegrationTypeResponse representa la respuesta de un tipo de integración
type IntegrationTypeResponse struct {
	ID                uint                         `json:"id" example:"1"`
	Name              string                       `json:"name" example:"WhatsApp"`
	Code              string                       `json:"code" example:"whatsapp"`
	Description       string                       `json:"description" example:"Integración con WhatsApp Cloud API"`
	Icon              string                       `json:"icon" example:"whatsapp-icon"`
	ImageURL          string                       `json:"image_url" example:"https://s3.amazonaws.com/bucket/integration-types/1234567890_logo.png"` // URL completa de la imagen
	Category          *IntegrationCategoryResponse `json:"category"`
	IsActive          bool                         `json:"is_active" example:"true"`
	ConfigSchema      datatypes.JSON               `json:"config_schema"`
	CredentialsSchema datatypes.JSON               `json:"credentials_schema"`
	SetupInstructions string                       `json:"setup_instructions" example:"1. Ve a Meta Business Suite\n2. Configura WhatsApp\n3. Copia credenciales"`
	CreatedAt         string                       `json:"created_at"`
	UpdatedAt         string                       `json:"updated_at"`
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
