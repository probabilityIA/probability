package response

import "time"

// IntegrationResponse representa la respuesta de una integración (sin credenciales)
type IntegrationResponse struct {
	ID          uint                   `json:"id" example:"1"`
	Name        string                 `json:"name" example:"WhatsApp Principal"`
	Code        string                 `json:"code" example:"whatsapp_platform"`
	Type        string                 `json:"type" example:"whatsapp"`
	Category    string                 `json:"category" example:"internal"`
	BusinessID  *uint                  `json:"business_id" example:"16"`
	IsActive    bool                   `json:"is_active" example:"true"`
	IsDefault   bool                   `json:"is_default" example:"true"`
	Config      map[string]interface{} `json:"config"`
	Description string                 `json:"description" example:"Integración principal de WhatsApp"`
	CreatedByID uint                   `json:"created_by_id" example:"1"`
	UpdatedByID *uint                  `json:"updated_by_id"`
	CreatedAt   time.Time              `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt   time.Time              `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// IntegrationListResponse representa la respuesta de lista de integraciones
type IntegrationListResponse struct {
	Success    bool                  `json:"success" example:"true"`
	Message    string                `json:"message" example:"Integraciones obtenidas exitosamente"`
	Data       []IntegrationResponse `json:"data"`
	Total      int64                 `json:"total" example:"25"`
	Page       int                   `json:"page" example:"1"`
	PageSize   int                   `json:"page_size" example:"10"`
	TotalPages int                   `json:"total_pages" example:"3"`
}

// IntegrationSuccessResponse representa la respuesta exitosa de una integración
type IntegrationSuccessResponse struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"Integración obtenida exitosamente"`
	Data    IntegrationResponse `json:"data"`
}

// IntegrationErrorResponse representa la respuesta de error
type IntegrationErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Error al procesar la solicitud"`
	Error   string `json:"error,omitempty" example:"Detalles del error"`
}

// IntegrationMessageResponse representa la respuesta de mensaje
type IntegrationMessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operación realizada exitosamente"`
}
