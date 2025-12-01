package response

// BusinessTokenResponse representa la respuesta del token de business
type BusinessTokenResponse struct {
	Token string `json:"token"`
}

// GenerateBusinessTokenSuccessResponse representa la respuesta exitosa
type GenerateBusinessTokenSuccessResponse struct {
	Success bool                  `json:"success"`
	Data    BusinessTokenResponse `json:"data"`
	Message string                `json:"message"`
}

// GenerateBusinessTokenErrorResponse representa la respuesta de error
type GenerateBusinessTokenErrorResponse struct {
	Error string `json:"error"`
}
