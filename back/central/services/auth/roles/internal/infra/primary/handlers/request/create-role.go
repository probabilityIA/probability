package request

// CreateRoleRequest representa la estructura para crear un nuevo rol
type CreateRoleRequest struct {
	Name           string `json:"name" binding:"required" example:"Administrador"`
	Description    string `json:"description" binding:"required" example:"Rol de administrador del sistema"`
	Level          int    `json:"level" binding:"required,min=1,max=10" example:"2"`
	IsSystem       bool   `json:"is_system" example:"false"`
	ScopeID        uint   `json:"scope_id" binding:"required" example:"1"`
	BusinessTypeID uint   `json:"business_type_id" binding:"required" example:"1"`
}
