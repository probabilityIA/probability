package response

import "time"

// BusinessResponse representa un negocio en la respuesta API
type BusinessResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Address         string    `json:"address"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Website         string    `json:"website"`
	LogoURL         string    `json:"logo_url"`
	PrimaryColor    string    `json:"primary_color"`
	SecondaryColor  string    `json:"secondary_color"`
	TertiaryColor   string    `json:"tertiary_color"`
	QuaternaryColor string    `json:"quaternary_color"`
	NavbarImageURL  string    `json:"navbar_image_url"`
	IsActive        bool      `json:"is_active"`
	BusinessTypeID  uint      `json:"business_type_id"`
	BusinessType    string    `json:"business_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// GetBusinessesResponse representa la respuesta para obtener múltiples negocios
type GetBusinessesResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
	Data       []BusinessResponse `json:"data"`
	Pagination *PaginationInfo    `json:"pagination,omitempty"`
}

// PaginationInfo representa información de paginación
type PaginationInfo struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// GetBusinessResponse representa la respuesta para obtener un negocio
type GetBusinessResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    BusinessResponse `json:"data"`
}

// CreateBusinessResponse representa la respuesta para crear un negocio
type CreateBusinessResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    BusinessResponse `json:"data"`
}

// UpdateBusinessResponse representa la respuesta para actualizar un negocio
type UpdateBusinessResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    BusinessResponse `json:"data"`
}

// DeleteBusinessResponse representa la respuesta para eliminar un negocio
type DeleteBusinessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
