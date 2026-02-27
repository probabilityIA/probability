package request

// CreateChannelStatusRequest representa los datos para crear un estado de canal de integración
type CreateChannelStatusRequest struct {
	IntegrationTypeID uint   `json:"integration_type_id" binding:"required"`
	Code              string `json:"code" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	IsActive          bool   `json:"is_active"`
	DisplayOrder      int    `json:"display_order"`
}

// UpdateChannelStatusRequest representa los datos para actualizar un estado de canal de integración
type UpdateChannelStatusRequest struct {
	Code         string `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	IsActive     bool   `json:"is_active"`
	DisplayOrder int    `json:"display_order"`
}
