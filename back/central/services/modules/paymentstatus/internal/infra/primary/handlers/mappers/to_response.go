package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/infra/primary/handlers/response"
)

// ToPaymentStatusResponse convierte DTO de dominio a respuesta HTTP
func ToPaymentStatusResponse(dto dtos.PaymentStatusInfo) response.PaymentStatusResponse {
	return response.PaymentStatusResponse{
		ID:          dto.ID,
		Code:        dto.Code,
		Name:        dto.Name,
		Description: dto.Description,
		Category:    dto.Category,
		Color:       dto.Color,
	}
}

// ToPaymentStatusListResponse convierte lista de DTOs a respuesta HTTP
func ToPaymentStatusListResponse(dtos []dtos.PaymentStatusInfo) []response.PaymentStatusResponse {
	result := make([]response.PaymentStatusResponse, len(dtos))
	for i, dto := range dtos {
		result[i] = ToPaymentStatusResponse(dto)
	}
	return result
}
