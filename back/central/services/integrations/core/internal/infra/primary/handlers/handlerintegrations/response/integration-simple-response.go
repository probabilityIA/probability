package response

// IntegrationSimpleResponse representa una integración en formato simplificado para dropdowns/selectores
type IntegrationSimpleResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`        // Código del tipo de integración (whatsapp, shopify, etc.)
	BusinessID *uint  `json:"business_id"` // Puede ser null para integraciones globales
	IsActive   bool   `json:"is_active"`
}

// GetIntegrationsSimpleResponse representa la respuesta para obtener integraciones en formato simple
type GetIntegrationsSimpleResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    []IntegrationSimpleResponse `json:"data"`
}
