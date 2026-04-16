package response

// AssignPermissionsToRoleResponse representa la respuesta al asignar permisos a un rol
type AssignPermissionsToRoleResponse struct {
	Success       bool   `json:"success" example:"true"`
	Message       string `json:"message" example:"Permisos asignados exitosamente al rol"`
	RoleID        uint   `json:"role_id" example:"1"`
	PermissionIDs []uint `json:"permission_ids"`
}

// GetRolePermissionsResponse representa la respuesta con los permisos de un rol
type GetRolePermissionsResponse struct {
	Success     bool                 `json:"success" example:"true"`
	Message     string               `json:"message" example:"Permisos del rol obtenidos exitosamente"`
	RoleID      uint                 `json:"role_id" example:"1"`
	RoleName    string               `json:"role_name" example:"Administrador"`
	Permissions []PermissionResponse `json:"permissions"`
	Count       int                  `json:"count" example:"3"`
}

// PermissionResponse representa un permiso en la respuesta
type PermissionResponse struct {
	ID          uint   `json:"id" example:"1"`
	Resource    string `json:"resource" example:"users"`
	Action      string `json:"action" example:"create"`
	Description string `json:"description" example:"Crear usuarios"`
	ScopeID     uint   `json:"scope_id" example:"1"`
	ScopeName   string `json:"scope_name" example:"Sistema"`
	ScopeCode   string `json:"scope_code" example:"system"`
}
