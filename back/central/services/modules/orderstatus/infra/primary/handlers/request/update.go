package request

// UpdateOrderStatusMappingRequest representa la solicitud para actualizar un mapeo de estado
type UpdateOrderStatusMappingRequest struct {
	OriginalStatus string `json:"original_status" binding:"required,max=128"`
	OrderStatusID  uint   `json:"order_status_id" binding:"required"`
	Priority       int    `json:"priority"`
	Description    string `json:"description"`
}
