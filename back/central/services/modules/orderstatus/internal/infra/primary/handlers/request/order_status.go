package request

// CreateOrderStatusRequest representa los datos para crear un estado de orden de Probability
type CreateOrderStatusRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
	Priority    int    `json:"priority"`
	IsActive    bool   `json:"is_active"`
}

// UpdateOrderStatusRequest representa los datos para actualizar un estado de orden de Probability
type UpdateOrderStatusRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
	Priority    int    `json:"priority"`
	IsActive    bool   `json:"is_active"`
}
