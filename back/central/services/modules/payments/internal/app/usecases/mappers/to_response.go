package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/entities"
)

// EntityToResponse convierte entidad de dominio a DTO de respuesta
func EntityToResponse(entity *entities.PaymentMethod) dtos.PaymentMethodResponse {
	return dtos.PaymentMethodResponse{
		ID:          entity.ID,
		Code:        entity.Code,
		Name:        entity.Name,
		Description: entity.Description,
		Category:    entity.Category,
		Provider:    entity.Provider,
		IsActive:    entity.IsActive,
		Icon:        entity.Icon,
		Color:       entity.Color,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// EntitiesToResponses convierte slice de entidades a slice de DTOs de respuesta
func EntitiesToResponses(entities []entities.PaymentMethod) []dtos.PaymentMethodResponse {
	responses := make([]dtos.PaymentMethodResponse, len(entities))
	for i, entity := range entities {
		responses[i] = EntityToResponse(&entity)
	}
	return responses
}

// MappingEntityToResponse convierte entidad de mapeo a DTO de respuesta
func MappingEntityToResponse(entity *entities.PaymentMethodMapping) dtos.PaymentMappingResponse {
	return dtos.PaymentMappingResponse{
		ID:              entity.ID,
		IntegrationType: entity.IntegrationType,
		OriginalMethod:  entity.OriginalMethod,
		PaymentMethodID: entity.PaymentMethodID,
		PaymentMethod:   EntityToResponse(&entity.PaymentMethod),
		IsActive:        entity.IsActive,
		Priority:        entity.Priority,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
	}
}

// MappingEntitiesToResponses convierte slice de entidades de mapeo a slice de DTOs
func MappingEntitiesToResponses(entities []entities.PaymentMethodMapping) []dtos.PaymentMappingResponse {
	responses := make([]dtos.PaymentMappingResponse, len(entities))
	for i, entity := range entities {
		responses[i] = MappingEntityToResponse(&entity)
	}
	return responses
}
