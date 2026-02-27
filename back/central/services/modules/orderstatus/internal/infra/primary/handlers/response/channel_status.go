package response

// IntegrationTypeResponse representa un tipo de integración ecommerce
type IntegrationTypeResponse struct {
	ID       uint   `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// IntegrationTypesResponse es el envelope de la lista de tipos de integración
type IntegrationTypesResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    []IntegrationTypeResponse `json:"data"`
}

// ChannelStatusResponse representa un estado nativo de un canal de integración
type ChannelStatusResponse struct {
	ID              uint                     `json:"id"`
	IntegrationTypeID uint                   `json:"integration_type_id"`
	IntegrationType   *IntegrationTypeResponse `json:"integration_type,omitempty"`
	Code            string                   `json:"code"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description,omitempty"`
	IsActive        bool                     `json:"is_active"`
	DisplayOrder    int                      `json:"display_order"`
}

// ChannelStatusListResponse es el envelope de la lista de estados de canal
type ChannelStatusListResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    []ChannelStatusResponse `json:"data"`
}

// ChannelStatusSingleResponse es el envelope de un solo estado de canal
type ChannelStatusSingleResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    ChannelStatusResponse `json:"data"`
}
