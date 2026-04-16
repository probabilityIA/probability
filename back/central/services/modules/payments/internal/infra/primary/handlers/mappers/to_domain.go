package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/request"
)

// RequestToCreatePaymentMethodDTO convierte request HTTP a DTO de dominio
func RequestToCreatePaymentMethodDTO(req *request.CreatePaymentMethod) *dtos.CreatePaymentMethod {
	return &dtos.CreatePaymentMethod{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Provider:    req.Provider,
		Icon:        req.Icon,
		Color:       req.Color,
	}
}

// RequestToUpdatePaymentMethodDTO convierte request HTTP a DTO de dominio
func RequestToUpdatePaymentMethodDTO(req *request.UpdatePaymentMethod) *dtos.UpdatePaymentMethod {
	return &dtos.UpdatePaymentMethod{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Provider:    req.Provider,
		Icon:        req.Icon,
		Color:       req.Color,
	}
}

// RequestToCreatePaymentMappingDTO convierte request HTTP a DTO de dominio
func RequestToCreatePaymentMappingDTO(req *request.CreatePaymentMapping) *dtos.CreatePaymentMapping {
	return &dtos.CreatePaymentMapping{
		IntegrationType: req.IntegrationType,
		OriginalMethod:  req.OriginalMethod,
		PaymentMethodID: req.PaymentMethodID,
		Priority:        req.Priority,
	}
}

// RequestToUpdatePaymentMappingDTO convierte request HTTP a DTO de dominio
func RequestToUpdatePaymentMappingDTO(req *request.UpdatePaymentMapping) *dtos.UpdatePaymentMapping {
	return &dtos.UpdatePaymentMapping{
		OriginalMethod:  req.OriginalMethod,
		PaymentMethodID: req.PaymentMethodID,
		Priority:        req.Priority,
	}
}
