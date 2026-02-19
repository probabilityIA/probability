package domain

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

type PermissionDTO struct {
	ID               uint
	Name             string
	Code             string
	Description      string
	Resource         string
	Action           string
	ResourceID       uint
	ActionID         uint
	ScopeID          uint
	ScopeName        string // Nombre del scope para mostrar
	ScopeCode        string // Código del scope para mostrar
	BusinessTypeID   uint   // ID del tipo de business
	BusinessTypeName string // Nombre del tipo de business
}

type CreatePermissionDTO struct {
	Name           string
	Code           string // Opcional, se genera automáticamente si no se proporciona
	Description    string
	ResourceID     uint // ID del resource
	ActionID       uint // ID de la action
	ScopeID        uint
	BusinessTypeID *uint // Opcional, nil = genérico
}

type BulkCreateResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type UpdatePermissionDTO struct {
	Name           string
	Code           string
	Description    string
	ResourceID     uint // ID del resource
	ActionID       uint // ID de la action
	ScopeID        uint
	BusinessTypeID *uint // Opcional, nil = genérico
}
