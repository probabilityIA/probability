package request

// GenerateAPIKeyRequest representa la solicitud para generar una API Key
type GenerateAPIKeyRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	BusinessID  uint   `json:"business_id" binding:"required"`
	Name        string `json:"name" binding:"required"` // Nombre de referencia de la API Key
	Description string `json:"description"`             // Descripción opcional
}
