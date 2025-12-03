package request

// CreateOrderStatusMappingRequest representa la solicitud para crear un mapeo de estado
type CreateOrderStatusMappingRequest struct {
	IntegrationType string `json:"integration_type" binding:"required,oneof=shopify whatsapp mercadolibre"`
	OriginalStatus  string `json:"original_status" binding:"required,max=128"`
	MappedStatus    string `json:"mapped_status" binding:"required,max=64"`
	Priority        int    `json:"priority"`
	Description     string `json:"description"`
}
