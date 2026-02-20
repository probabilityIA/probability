package request

// BulkCreatePermissionItem representa un permiso individual en una creación masiva
type BulkCreatePermissionItem struct {
	Name           string `json:"name" binding:"required"`
	ResourceID     uint   `json:"resource_id" binding:"required"`
	ActionID       uint   `json:"action_id" binding:"required"`
	ScopeID        uint   `json:"scope_id" binding:"required"`
	BusinessTypeID *uint  `json:"business_type_id"`
}

// BulkCreatePermissionRequest representa la solicitud de creación masiva de permisos
type BulkCreatePermissionRequest struct {
	Permissions []BulkCreatePermissionItem `json:"permissions" binding:"required,min=1"`
}
