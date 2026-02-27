package request

// CreateMovementTypeRequest payload para crear un tipo de movimiento
type CreateMovementTypeRequest struct {
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Direction   string `json:"direction" binding:"required,oneof=in out neutral"`
}

// UpdateMovementTypeRequest payload para actualizar un tipo de movimiento
type UpdateMovementTypeRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	IsActive    *bool  `json:"is_active"`
	Direction   string `json:"direction" binding:"omitempty,oneof=in out neutral"`
}
