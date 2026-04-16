package response

// ValidateAPIKeyResponse representa la respuesta de validación de API Key
type ValidateAPIKeyResponse struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	UserID     uint     `json:"user_id"`
	Email      string   `json:"email"`
	BusinessID uint     `json:"business_id"`
	Roles      []string `json:"roles"`
	APIKeyID   uint     `json:"api_key_id"`
}

// ValidateAPIKeySuccessResponse representa la respuesta exitosa para Swagger
type ValidateAPIKeySuccessResponse struct {
	Success bool                   `json:"success"`
	Data    ValidateAPIKeyResponse `json:"data"`
}

// ValidateAPIKeyErrorResponse representa la respuesta de error para Swagger
type ValidateAPIKeyErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
