package response

// RoleResponse representa la respuesta de un rol
type RoleResponse struct {
	ID               uint   `json:"id" example:"1"`
	Name             string `json:"name" example:"Administrador"`
	Code             string `json:"code" example:"admin"`
	Description      string `json:"description" example:"Rol de administrador del sistema"`
	Level            int    `json:"level" example:"2"`
	IsSystem         bool   `json:"is_system" example:"true"`
	ScopeID          uint   `json:"scope_id" example:"1"`
	ScopeName        string `json:"scope_name" example:"Sistema"`
	ScopeCode        string `json:"scope_code" example:"system"`
	BusinessTypeID   uint   `json:"business_type_id" example:"1"`
	BusinessTypeName string `json:"business_type_name" example:"Propiedad Horizontal"`
}

// RoleListResponse representa la respuesta de una lista de roles
type RoleListResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    []RoleResponse `json:"data"`
	Count   int            `json:"count" example:"5"`
}

// RoleSuccessResponse representa la respuesta exitosa de un rol individual
type RoleSuccessResponse struct {
	Success bool         `json:"success" example:"true"`
	Data    RoleResponse `json:"data"`
}

// RoleErrorResponse representa la respuesta de error para roles
type RoleErrorResponse struct {
	Error string `json:"error" example:"Error interno del servidor"`
}
