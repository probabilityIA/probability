package dtos

// ListWarehousesParams parámetros para listar bodegas
type ListWarehousesParams struct {
	BusinessID    uint
	IsActive      *bool
	IsFulfillment *bool
	Search        string
	Page          int
	PageSize      int
}

// Offset calcula el offset para paginación
func (p ListWarehousesParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// CreateWarehouseDTO datos para crear una bodega
type CreateWarehouseDTO struct {
	BusinessID    uint
	Name          string
	Code          string
	Address       string
	City          string
	State         string
	Country       string
	ZipCode       string
	Phone         string
	ContactName   string
	ContactEmail  string
	IsActive      bool
	IsDefault     bool
	IsFulfillment bool
}

// UpdateWarehouseDTO datos para actualizar una bodega
type UpdateWarehouseDTO struct {
	ID            uint
	BusinessID    uint
	Name          string
	Code          string
	Address       string
	City          string
	State         string
	Country       string
	ZipCode       string
	Phone         string
	ContactName   string
	ContactEmail  string
	IsActive      bool
	IsDefault     bool
	IsFulfillment bool
}
