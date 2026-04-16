package request

// CreatePermissionRequest representa la solicitud para crear un permiso
type CreatePermissionRequest struct {
	Name           string `json:"name" binding:"required" example:"Crear usuario"`
	Code           string `json:"code" example:"horizontalproperty_createuser"` // Opcional, se genera automáticamente si no se proporciona
	Description    string `json:"description" example:"Permite crear nuevos usuarios en el sistema"`
	ResourceID     uint   `json:"resource_id" binding:"required" example:"1"` // ID del resource
	ActionID       uint   `json:"action_id" binding:"required" example:"1"`   // ID de la action
	ScopeID        uint   `json:"scope_id" binding:"required" example:"1"`
	BusinessTypeID *uint  `json:"business_type_id" example:"11"`
}

// UpdatePermissionRequest representa la solicitud para actualizar un permiso
type UpdatePermissionRequest struct {
	Name           string `json:"name" binding:"required" example:"Crear usuarios"`
	Code           string `json:"code" example:"users:create"`
	Description    string `json:"description" example:"Permite crear nuevos usuarios en el sistema"`
	ResourceID     uint   `json:"resource_id" binding:"required" example:"1"` // ID del resource
	ActionID       uint   `json:"action_id" binding:"required" example:"1"`   // ID de la action
	ScopeID        uint   `json:"scope_id" binding:"required" example:"1"`
	BusinessTypeID *uint  `json:"business_type_id" example:"11"`
}
