package response

// Error representa la respuesta HTTP de error estándar
// ✅ DTO HTTP - CON TAGS (json)
type Error struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Success representa la respuesta HTTP de éxito con data
// ✅ DTO HTTP - CON TAGS (json)
type Success struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
