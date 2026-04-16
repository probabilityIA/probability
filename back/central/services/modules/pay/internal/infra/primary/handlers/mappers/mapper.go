package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/response"
)

// ToResponse convierte una entidad de dominio a respuesta HTTP
func ToResponse(tx *entities.PaymentTransaction) *response.PaymentTransactionResponse {
	return &response.PaymentTransactionResponse{
		ID:            tx.ID,
		BusinessID:    tx.BusinessID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		Status:        string(tx.Status),
		GatewayCode:   tx.GatewayCode,
		ExternalID:    tx.ExternalID,
		Reference:     tx.Reference,
		PaymentMethod: tx.PaymentMethod,
		Description:   tx.Description,
		CreatedAt:     tx.CreatedAt,
		UpdatedAt:     tx.UpdatedAt,
	}
}
