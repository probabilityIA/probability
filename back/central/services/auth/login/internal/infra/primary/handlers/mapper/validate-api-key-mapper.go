package mapper

import (
	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
)

// ToValidateAPIKeyResponse convierte el DTO del dominio a la respuesta HTTP
func ToValidateAPIKeyResponse(dto *domain.ValidateAPIKeyResponse) response.ValidateAPIKeyResponse {
	return response.ValidateAPIKeyResponse{
		Success:    dto.Success,
		Message:    dto.Message,
		UserID:     dto.UserID,
		Email:      dto.Email,
		BusinessID: dto.BusinessID,
		Roles:      dto.Roles,
		APIKeyID:   dto.APIKeyID,
	}
}

// ToValidateAPIKeySuccessResponse convierte el DTO del dominio a la respuesta exitosa HTTP
func ToValidateAPIKeySuccessResponse(dto *domain.ValidateAPIKeyResponse) response.ValidateAPIKeySuccessResponse {
	return response.ValidateAPIKeySuccessResponse{
		Success: true,
		Data:    ToValidateAPIKeyResponse(dto),
	}
}

// ToValidateAPIKeyErrorResponse crea una respuesta de error HTTP
func ToValidateAPIKeyErrorResponse(message string) response.ValidateAPIKeyErrorResponse {
	return response.ValidateAPIKeyErrorResponse{
		Success: false,
		Message: message,
	}
}
