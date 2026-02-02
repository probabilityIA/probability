package response

// IntegrationSimpleResponse representa una integración en formato simplificado para dropdowns/selectores
type IntegrationSimpleResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`          // Código del tipo de integración (whatsapp, shopify, etc.)
	Category      string `json:"category"`      // Código de la categoría (ecommerce, messaging, etc.)
	CategoryName  string `json:"category_name"` // Nombre de la categoría (E-commerce, Mensajería, etc.)
	CategoryColor string `json:"category_color,omitempty"` // Color hexadecimal de la categoría
	BusinessID    *uint  `json:"business_id"` // Puede ser null para integraciones globales
	IsActive      bool   `json:"is_active"`
}

// GetIntegrationsSimpleResponse representa la respuesta para obtener integraciones en formato simple
type GetIntegrationsSimpleResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    []IntegrationSimpleResponse `json:"data"`
}
