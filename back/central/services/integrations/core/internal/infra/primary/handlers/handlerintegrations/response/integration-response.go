package response

import "time"

// IntegrationTypeInfo representa información básica del tipo de integración
type IntegrationTypeInfo struct {
	ID       uint   `json:"id" example:"1"`
	Name     string `json:"name" example:"WhatsApp"`
	Code     string `json:"code" example:"whatsapp"`
	ImageURL string `json:"image_url" example:"https://s3.amazonaws.com/bucket/integration-types/1234567890_logo.png"` // URL completa de la imagen
}

// IntegrationResponse representa la respuesta de una integración (sin credenciales)
type IntegrationResponse struct {
	ID                uint                   `json:"id" example:"1"`
	Name              string                 `json:"name" example:"WhatsApp Principal"`
	Code              string                 `json:"code" example:"whatsapp_platform"`
	IntegrationTypeID uint                   `json:"integration_type_id" example:"1"`
	IntegrationType   *IntegrationTypeInfo   `json:"integration_type,omitempty"` // Información del tipo si está cargado
	Category          string                 `json:"category" example:"internal"`
	BusinessID        *uint                  `json:"business_id" example:"16"`
	StoreID           string                 `json:"store_id" example:"mystore.myshopify.com"` // Identificador externo (ej: shop domain)
	IsActive          bool                   `json:"is_active" example:"true"`
	IsDefault         bool                   `json:"is_default" example:"true"`
	Config            map[string]interface{} `json:"config"`
	Credentials       map[string]interface{} `json:"credentials,omitempty"` // Solo se incluye cuando se solicita para edición
	Description       string                 `json:"description" example:"Integración principal de WhatsApp"`
	CreatedByID       uint                   `json:"created_by_id" example:"1"`
	UpdatedByID       *uint                  `json:"updated_by_id"`
	CreatedAt         time.Time              `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt         time.Time              `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// IntegrationListResponse representa la respuesta de lista de integraciones
//
//	@Description	Respuesta de lista paginada de integraciones
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
//
//	@Description	Respuesta exitosa con datos de integración
type IntegrationSuccessResponse struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"Integración obtenida exitosamente"`
	Data    IntegrationResponse `json:"data"`
}

// IntegrationErrorResponse representa la respuesta de error
//
//	@Description	Respuesta de error para operaciones fallidas
type IntegrationErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Error al procesar la solicitud"`
	Error   string `json:"error,omitempty" example:"Detalles del error"`
}

// IntegrationMessageResponse representa la respuesta de mensaje
//
//	@Description	Respuesta de mensaje para operaciones simples
type IntegrationMessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operación realizada exitosamente"`
}

// ErrorResponse representa la respuesta de error genérica
type ErrorResponse struct {
	Error string `json:"error" example:"Error al procesar la solicitud"`
}

// WebhookURLData contiene la información del webhook
type WebhookURLData struct {
	URL         string   `json:"url" example:"https://api.example.com/integrations/shopify/webhook"`
	Method      string   `json:"method" example:"POST"`
	Description string   `json:"description" example:"URL para configurar en Shopify para recibir eventos de órdenes"`
	Events      []string `json:"events,omitempty" example:"orders/create,orders/updated"`
}

// WebhookURLResponse representa la respuesta con la URL del webhook
//
//	@Description	Respuesta con la información del webhook para configurar en la plataforma externa
type WebhookURLResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    *WebhookURLData `json:"data"`
}

// ListWebhooksResponse representa la respuesta con la lista de webhooks
//
//	@Description	Respuesta con la lista de webhooks configurados para una integración
type ListWebhooksResponse struct {
	Success bool          `json:"success" example:"true"`
	Data    []interface{} `json:"data"`
}

// DeleteWebhookResponse representa la respuesta al eliminar un webhook
//
//	@Description	Respuesta de éxito al eliminar un webhook
type DeleteWebhookResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Webhook eliminado exitosamente"`
}

// VerifyWebhooksResponse representa la respuesta al verificar webhooks existentes
//
//	@Description	Respuesta con la lista de webhooks que coinciden con nuestra URL
type VerifyWebhooksResponse struct {
	Success bool          `json:"success" example:"true"`
	Data    []interface{} `json:"data"` // Lista de webhooks que coinciden
	Message string        `json:"message" example:"Webhooks verificados exitosamente"`
}

// CreateWebhookResponseData contiene los datos del resultado de crear webhooks
type CreateWebhookResponseData struct {
	ExistingWebhooks []interface{} `json:"existing_webhooks"` // Webhooks encontrados que coinciden
	DeletedWebhooks  []interface{} `json:"deleted_webhooks"`  // Webhooks eliminados
	CreatedWebhooks  []string      `json:"created_webhooks"`  // IDs de webhooks creados
	WebhookURL       string        `json:"webhook_url"`       // URL del webhook
}

// CreateWebhookResponse representa la respuesta al crear webhooks
//
//	@Description	Respuesta con información sobre webhooks encontrados, eliminados y creados
type CreateWebhookResponse struct {
	Success bool                      `json:"success" example:"true"`
	Data    CreateWebhookResponseData `json:"data"`
	Message string                    `json:"message" example:"Webhooks creados exitosamente"`
}
