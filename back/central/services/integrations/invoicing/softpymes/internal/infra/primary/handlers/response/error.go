package response

// Error es la respuesta de error estándar
type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
