package dtos

// ListClientsParams parámetros de búsqueda y paginación para listar clientes
type ListClientsParams struct {
	BusinessID uint
	Search     string
	Page       int
	PageSize   int
}

// Offset calcula el offset para paginación
func (p ListClientsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// CreateClientDTO datos para crear un cliente
type CreateClientDTO struct {
	BusinessID uint
	Name       string
	Email      string
	Phone      string
	Dni        *string
}

// UpdateClientDTO datos para actualizar un cliente
type UpdateClientDTO struct {
	ID         uint
	BusinessID uint
	Name       string
	Email      string
	Phone      string
	Dni        *string
}
