package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/infra/primary/handlers/response"
)

// DomainToPaymentMethodResponse convierte DTO de dominio a respuesta HTTP
func DomainToPaymentMethodResponse(dto *dtos.PaymentMethodResponse) *response.PaymentMethod {
	return &response.PaymentMethod{
		ID:          dto.ID,
		Code:        dto.Code,
		Name:        dto.Name,
		Description: dto.Description,
		Category:    dto.Category,
		Provider:    dto.Provider,
		IsActive:    dto.IsActive,
		Icon:        dto.Icon,
		Color:       dto.Color,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

// DomainToPaymentMethodsListResponse convierte DTO de dominio a respuesta HTTP paginada
func DomainToPaymentMethodsListResponse(dto *dtos.PaymentMethodsListResponse) *response.PaymentMethodsList {
	data := make([]response.PaymentMethod, len(dto.Data))
	for i, item := range dto.Data {
		data[i] = *DomainToPaymentMethodResponse(&item)
	}

	return &response.PaymentMethodsList{
		Data:       data,
		Total:      dto.Total,
		Page:       dto.Page,
		PageSize:   dto.PageSize,
		TotalPages: dto.TotalPages,
	}
}

// DomainToPaymentMappingResponse convierte DTO de dominio a respuesta HTTP
func DomainToPaymentMappingResponse(dto *dtos.PaymentMappingResponse) *response.PaymentMapping {
	return &response.PaymentMapping{
		ID:              dto.ID,
		IntegrationType: dto.IntegrationType,
		OriginalMethod:  dto.OriginalMethod,
		PaymentMethodID: dto.PaymentMethodID,
		PaymentMethod:   *DomainToPaymentMethodResponse(&dto.PaymentMethod),
		IsActive:        dto.IsActive,
		Priority:        dto.Priority,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}

// DomainToPaymentMappingsListResponse convierte DTO de dominio a respuesta HTTP de lista
func DomainToPaymentMappingsListResponse(dto *dtos.PaymentMappingsListResponse) *response.PaymentMappingsList {
	data := make([]response.PaymentMapping, len(dto.Data))
	for i, item := range dto.Data {
		data[i] = *DomainToPaymentMappingResponse(&item)
	}

	return &response.PaymentMappingsList{
		Data:  data,
		Total: dto.Total,
	}
}

// DomainToPaymentMappingsByIntegrationResponse convierte slice de DTOs agrupados a respuesta HTTP
func DomainToPaymentMappingsByIntegrationResponse(dtos []dtos.PaymentMappingsByIntegrationResponse) []response.PaymentMappingsByIntegration {
	result := make([]response.PaymentMappingsByIntegration, len(dtos))
	for i, dto := range dtos {
		mappings := make([]response.PaymentMapping, len(dto.Mappings))
		for j, mapping := range dto.Mappings {
			mappings[j] = *DomainToPaymentMappingResponse(&mapping)
		}

		result[i] = response.PaymentMappingsByIntegration{
			IntegrationType: dto.IntegrationType,
			Mappings:        mappings,
		}
	}
	return result
}
