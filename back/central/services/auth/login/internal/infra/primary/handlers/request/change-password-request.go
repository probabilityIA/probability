package request

// ChangePasswordRequest representa la solicitud para cambiar contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=6,max=100"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=100"`
}
