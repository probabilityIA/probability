package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/entities"
)

// CreateDTOToEntity convierte DTO de creación a entidad de dominio
func CreateDTOToEntity(dto *dtos.CreatePaymentMethod) *entities.PaymentMethod {
	return &entities.PaymentMethod{
		Code:        dto.Code,
		Name:        dto.Name,
		Description: dto.Description,
		Category:    dto.Category,
		Provider:    dto.Provider,
		Icon:        dto.Icon,
		Color:       dto.Color,
		IsActive:    true,
	}
}

// CreateMappingDTOToEntity convierte DTO de creación de mapeo a entidad
func CreateMappingDTOToEntity(dto *dtos.CreatePaymentMapping) *entities.PaymentMethodMapping {
	return &entities.PaymentMethodMapping{
		IntegrationType: dto.IntegrationType,
		OriginalMethod:  dto.OriginalMethod,
		PaymentMethodID: dto.PaymentMethodID,
		Priority:        dto.Priority,
		IsActive:        true,
	}
}
