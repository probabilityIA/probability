package domain

import "time"

type Role struct {
	ID               uint
	Name             string
	Description      string
	Level            int
	IsSystem         bool
	ScopeID          uint
	ScopeName        string
	ScopeCode        string
	BusinessTypeID   uint
	BusinessTypeName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type RoleDTO struct {
	ID               uint
	Name             string
	Code             string
	Description      string
	Level            int
	IsSystem         bool
	ScopeID          uint
	ScopeName        string // Nombre del scope para mostrar
	ScopeCode        string // Código del scope para mostrar
	BusinessTypeID   uint   // ID del tipo de business
	BusinessTypeName string // Nombre del tipo de business
}

type CreateRoleDTO struct {
	Name           string
	Description    string
	Level          int
	IsSystem       bool
	ScopeID        uint
	BusinessTypeID uint
}

type UpdateRoleDTO struct {
	Name           *string
	Description    *string
	Level          *int
	IsSystem       *bool
	ScopeID        *uint
	BusinessTypeID *uint
}

type RoleFilters struct {
	Name           *string
	Level          *int
	IsSystem       *bool
	ScopeID        *uint
	BusinessTypeID *uint
}

type Permission struct {
	ID               uint
	Name             string
	Code             string
	Description      string
	Resource         string
	Action           string
	ResourceID       uint
	ActionID         uint // ID de la acción
	ScopeID          uint
	ScopeName        string
	ScopeCode        string
	BusinessTypeID   uint
	BusinessTypeName string
}
