package mapper

import (
	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
)

// ToGenerateAPIKeyRequest convierte el request HTTP a DTO de dominio
func ToGenerateAPIKeyRequest(req request.GenerateAPIKeyRequest, requesterID uint) domain.GenerateAPIKeyRequest {
	return domain.GenerateAPIKeyRequest{
		UserID:      req.UserID,
		BusinessID:  req.BusinessID,
		Name:        req.Name,
		Description: req.Description,
		RequesterID: requesterID,
	}
}

// ToGenerateAPIKeyResponse convierte el DTO de dominio a response HTTP
func ToGenerateAPIKeyResponse(dto *domain.GenerateAPIKeyResponse) response.GenerateAPIKeyResponse {
	return response.GenerateAPIKeyResponse{
		Success:     dto.Success,
		Message:     dto.Message,
		APIKey:      dto.APIKey,
		UserID:      dto.APIKeyInfo.UserID,
		BusinessID:  dto.APIKeyInfo.BusinessID,
		Name:        dto.APIKeyInfo.Name,
		Description: dto.APIKeyInfo.Description,
		RateLimit:   dto.APIKeyInfo.RateLimit,
		CreatedAt:   dto.APIKeyInfo.CreatedAt,
	}
}
