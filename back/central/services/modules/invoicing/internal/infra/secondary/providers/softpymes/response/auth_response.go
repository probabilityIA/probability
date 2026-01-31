package response

// AuthResponse representa la respuesta de autenticación de Softpymes
type AuthResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	ExpiresIn int  `json:"expires_in"` // Tiempo de expiración en segundos
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
