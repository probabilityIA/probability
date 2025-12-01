package response

// ChangePasswordResponse representa la respuesta al cambiar contraseña
type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
