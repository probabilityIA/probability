package dtos

// ListFilters representa filtros para listar mapeos
// PURO - Sin tags de validación
type ListFilters struct {
	Page              int
	PageSize          int
	IntegrationTypeID *uint
	IsActive          *bool
}
