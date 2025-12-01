package request

// GeneratePasswordRequest representa la solicitud para generar una nueva contraseña aleatoria
// Si el usuario es super admin, puede especificar user_id para generar contraseña de otro usuario
// Si no se envía user_id, se genera la contraseña para el usuario autenticado
type GeneratePasswordRequest struct {
	UserID *uint `json:"user_id" binding:"omitempty,min=1" example:"10"` // Opcional: solo para super usuarios
}
