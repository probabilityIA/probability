package domain

import "time"

// Resource representa un recurso del sistema
type Resource struct {
	ID               uint
	Name             string
	Description      string
	BusinessTypeID   uint
	BusinessTypeName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
