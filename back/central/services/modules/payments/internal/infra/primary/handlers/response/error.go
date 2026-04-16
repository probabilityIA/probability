package response

// Error representa una respuesta HTTP de error
type Error struct {
	Error string `json:"error"`
}
