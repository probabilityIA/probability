package request

// CreateClientRequest payload de creación de cliente
type CreateClientRequest struct {
	Name  string  `json:"name" binding:"required,min=2,max=255"`
	Email *string `json:"email" binding:"omitempty,email,max=255"`
	Phone string  `json:"phone" binding:"omitempty,max=20"`
	Dni   *string `json:"dni" binding:"omitempty,max=30"`
}

// UpdateClientRequest payload de actualización de cliente
type UpdateClientRequest struct {
	Name  string  `json:"name" binding:"required,min=2,max=255"`
	Email *string `json:"email" binding:"omitempty,email,max=255"`
	Phone string  `json:"phone" binding:"omitempty,max=20"`
	Dni   *string `json:"dni" binding:"omitempty,max=30"`
}
