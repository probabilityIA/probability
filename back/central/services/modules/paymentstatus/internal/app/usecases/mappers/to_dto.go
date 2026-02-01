package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/entities"
)

// ToPaymentStatusInfo convierte una entidad a DTO de información
func ToPaymentStatusInfo(entity entities.PaymentStatus) dtos.PaymentStatusInfo {
	return dtos.PaymentStatusInfo{
		ID:          entity.ID,
		Code:        entity.Code,
		Name:        entity.Name,
		Description: entity.Description,
		Category:    entity.Category,
		Color:       entity.Color,
	}
}

// ToPaymentStatusInfoList convierte lista de entidades a DTOs
func ToPaymentStatusInfoList(entities []entities.PaymentStatus) []dtos.PaymentStatusInfo {
	result := make([]dtos.PaymentStatusInfo, len(entities))
	for i, entity := range entities {
		result[i] = ToPaymentStatusInfo(entity)
	}
	return result
}
