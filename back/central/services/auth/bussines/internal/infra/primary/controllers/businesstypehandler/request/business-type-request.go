package request

// BusinessTypeRequest representa la solicitud para crear/actualizar un tipo de negocio
type BusinessTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code"` // Opcional, se genera automáticamente si no se proporciona
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IsActive    bool   `json:"is_active"`
}
