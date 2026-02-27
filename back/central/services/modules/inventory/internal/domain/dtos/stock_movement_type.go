package dtos

// CreateStockMovementTypeDTO datos para crear un tipo de movimiento
type CreateStockMovementTypeDTO struct {
	Code        string
	Name        string
	Description string
	Direction   string // in, out, neutral
}

// UpdateStockMovementTypeDTO datos para actualizar un tipo de movimiento
type UpdateStockMovementTypeDTO struct {
	ID          uint
	Name        string
	Description string
	IsActive    *bool
	Direction   string
}

// ListStockMovementTypesParams parámetros para listar tipos de movimiento
type ListStockMovementTypesParams struct {
	BusinessID uint // no se usa en filtro (los tipos son globales), pero se requiere para auth
	ActiveOnly bool
	Page       int
	PageSize   int
}

// Offset calcula el offset para paginación
func (p ListStockMovementTypesParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}
