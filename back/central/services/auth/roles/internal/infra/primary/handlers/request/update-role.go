package request

// UpdateRoleRequest representa la estructura para actualizar un rol existente
type UpdateRoleRequest struct {
	Name           *string `json:"name" example:"Administrador Actualizado"`
	Description    *string `json:"description" example:"Rol de administrador actualizado"`
	Level          *int    `json:"level" binding:"omitempty,min=1,max=10" example:"3"`
	IsSystem       *bool   `json:"is_system" example:"false"`
	ScopeID        *uint   `json:"scope_id" example:"1"`
	BusinessTypeID *uint   `json:"business_type_id" example:"1"`
}

