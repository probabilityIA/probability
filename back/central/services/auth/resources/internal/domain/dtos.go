package domain

import "time"

// ResourceDTO representa un recurso para casos de uso
type ResourceDTO struct {
	ID               uint
	Name             string
	Description      string
	BusinessTypeID   uint
	BusinessTypeName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CreateResourceDTO representa los datos para crear un recurso
type CreateResourceDTO struct {
	Name           string
	Description    string
	BusinessTypeID *uint // Opcional, nil = genérico
}

// UpdateResourceDTO representa los datos para actualizar un recurso
type UpdateResourceDTO struct {
	Name           string
	Description    string
	BusinessTypeID *uint // Opcional, nil = genérico
}

// ResourceListDTO representa una lista paginada de recursos
type ResourceListDTO struct {
	Resources  []ResourceDTO
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ResourceFilters representa los filtros para la consulta de recursos
type ResourceFilters struct {
	Page           int
	PageSize       int
	Name           string
	Description    string
	BusinessTypeID *uint  // Filtrar por tipo de business (null = genérico, aplica a todos)
	SortBy         string // "name", "created_at", etc.
	SortOrder      string // "asc" o "desc"
}
