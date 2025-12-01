package response

import "time"

// BusinessTypeResponse representa un tipo de negocio en la respuesta API
type BusinessTypeResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetBusinessTypesResponse representa la respuesta para obtener múltiples tipos de negocio
type GetBusinessTypesResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    []BusinessTypeResponse `json:"data"`
}

// GetBusinessTypeResponse representa la respuesta para obtener un tipo de negocio
type GetBusinessTypeResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    BusinessTypeResponse `json:"data"`
}

// CreateBusinessTypeResponse representa la respuesta para crear un tipo de negocio
type CreateBusinessTypeResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    BusinessTypeResponse `json:"data"`
}

// UpdateBusinessTypeResponse representa la respuesta para actualizar un tipo de negocio
type UpdateBusinessTypeResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    BusinessTypeResponse `json:"data"`
}

// DeleteBusinessTypeResponse representa la respuesta para eliminar un tipo de negocio
type DeleteBusinessTypeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
