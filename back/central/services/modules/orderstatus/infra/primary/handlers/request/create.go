package request

// CreateOrderStatusMappingRequest representa la solicitud para crear un mapeo de estado
type CreateOrderStatusMappingRequest struct {
	IntegrationTypeID uint   `json:"integration_type_id" binding:"required"`
	OriginalStatus    string `json:"original_status" binding:"required,max=128"`
	OrderStatusID     uint   `json:"order_status_id" binding:"required"`
	Priority          int    `json:"priority"`
	Description       string `json:"description"`
}
