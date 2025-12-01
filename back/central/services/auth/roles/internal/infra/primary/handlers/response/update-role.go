package response

// UpdateRoleResponse representa la respuesta al actualizar un rol
type UpdateRoleResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"Rol actualizado exitosamente"`
	Data    RoleData `json:"data"`
}
