package request

// AssignPermissionsToRoleRequest representa la solicitud para asignar permisos a un rol
type AssignPermissionsToRoleRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required,dive,min=1"`
}
