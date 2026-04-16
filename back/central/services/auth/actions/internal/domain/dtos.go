package domain

import "time"

// ActionDTO representa un action para casos de uso
type ActionDTO struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateActionDTO representa los datos para crear un action
type CreateActionDTO struct {
	Name        string
	Description string
}

// UpdateActionDTO representa los datos para actualizar un action
type UpdateActionDTO struct {
	Name        string
	Description string
}

// ActionListDTO representa una lista paginada de actions
type ActionListDTO struct {
	Actions    []ActionDTO
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
