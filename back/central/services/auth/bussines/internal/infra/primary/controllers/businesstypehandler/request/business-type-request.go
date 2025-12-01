package request

// BusinessTypeRequest representa la solicitud para crear/actualizar un tipo de negocio
type BusinessTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IsActive    bool   `json:"is_active"`
}
