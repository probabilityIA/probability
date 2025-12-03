package request

// UpdateOrderStatusMappingRequest representa la solicitud para actualizar un mapeo de estado
type UpdateOrderStatusMappingRequest struct {
	OriginalStatus string `json:"original_status" binding:"required,max=128"`
	MappedStatus   string `json:"mapped_status" binding:"required,max=64"`
	Priority       int    `json:"priority"`
	Description    string `json:"description"`
}
