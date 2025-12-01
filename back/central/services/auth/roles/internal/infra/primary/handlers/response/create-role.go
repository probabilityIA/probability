package response

import "time"

// CreateRoleResponse representa la respuesta al crear un rol
type CreateRoleResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Rol creado exitosamente"`
	Data    RoleData `json:"data"`
}

// RoleData contiene los datos del rol creado
type RoleData struct {
	ID             uint      `json:"id" example:"1"`
	Name           string    `json:"name" example:"Administrador"`
	Description    string    `json:"description" example:"Rol de administrador del sistema"`
	Level          int       `json:"level" example:"2"`
	IsSystem       bool      `json:"is_system" example:"false"`
	ScopeID        uint      `json:"scope_id" example:"1"`
	BusinessTypeID uint      `json:"business_type_id" example:"1"`
	CreatedAt      time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
