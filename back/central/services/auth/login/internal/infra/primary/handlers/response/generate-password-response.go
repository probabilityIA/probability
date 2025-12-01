package response

// GeneratePasswordResponse representa la respuesta al generar una nueva contraseña aleatoria
// La contraseña solo se muestra una vez en esta respuesta
type GeneratePasswordResponse struct {
	Success  bool   `json:"success" example:"true"`
	Email    string `json:"email" example:"usuario@ejemplo.com"`
	Password string `json:"password" example:"aB3$kL9mP2xQ"`
	Message  string `json:"message" example:"Nueva contraseña generada para el usuario usuario@ejemplo.com"`
}
