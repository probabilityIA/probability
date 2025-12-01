package response

import "time"

// GenerateAPIKeyResponse representa la respuesta de generación de API Key
type GenerateAPIKeyResponse struct {
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
	APIKey      string    `json:"api_key"`
	UserID      uint      `json:"user_id"`
	BusinessID  uint      `json:"business_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RateLimit   int       `json:"rate_limit"`
	CreatedAt   time.Time `json:"created_at"`
}

// GenerateAPIKeySuccessResponse representa la respuesta exitosa para Swagger
type GenerateAPIKeySuccessResponse struct {
	Success bool                   `json:"success"`
	Data    GenerateAPIKeyResponse `json:"data"`
}

// GenerateAPIKeyErrorResponse representa la respuesta de error para Swagger
type GenerateAPIKeyErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
