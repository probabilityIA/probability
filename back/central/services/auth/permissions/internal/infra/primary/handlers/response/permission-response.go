package response

// PermissionResponse representa la respuesta de un permiso
type PermissionResponse struct {
	ID               uint   `json:"id" example:"1"`
	Name             string `json:"name" example:"Crear usuarios"`
	Code             string `json:"code" example:"users:create"`
	Description      string `json:"description" example:"Permite crear nuevos usuarios en el sistema"`
	Resource         string `json:"resource" example:"users"`
	Action           string `json:"action" example:"create"`
	ResourceID       uint   `json:"resource_id" example:"3"`
	ActionID         uint   `json:"action_id" example:"5"`
	ScopeID          uint   `json:"scope_id" example:"1"`
	ScopeName        string `json:"scope_name" example:"Sistema"`
	ScopeCode        string `json:"scope_code" example:"system"`
	BusinessTypeID   uint   `json:"business_type_id" example:"11"`
	BusinessTypeName string `json:"business_type_name" example:"Propiedad Horizontal"`
}

// PermissionListResponse representa la respuesta de una lista de permisos
type PermissionListResponse struct {
	Success bool                 `json:"success" example:"true"`
	Data    []PermissionResponse `json:"data"`
	Total   int                  `json:"total" example:"10"`
}

// PermissionSuccessResponse representa la respuesta exitosa de un permiso
type PermissionSuccessResponse struct {
	Success bool               `json:"success" example:"true"`
	Data    PermissionResponse `json:"data"`
}

// PermissionErrorResponse representa la respuesta de error
type PermissionErrorResponse struct {
	Error string `json:"error" example:"Error al procesar la solicitud"`
}

// PermissionMessageResponse representa la respuesta de mensaje para operaciones CRUD
type PermissionMessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Permiso creado exitosamente"`
}
