package response

// AssignRoleToUserBusinessResponse representa la respuesta al asignar un rol a un usuario en un business
type AssignRoleToUserBusinessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Rol asignado exitosamente al usuario en el business"`
}
