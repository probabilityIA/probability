package mapper

import (
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/request"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/response"
)

// RequestToDTO convierte request.BusinessTypeRequest a dtos.BusinessTypeRequest
func RequestToDTO(req request.BusinessTypeRequest) domain.BusinessTypeRequest {
	return domain.BusinessTypeRequest{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
	}
}

// BusinessTypeDTOToResponse convierte un DTO BusinessTypeResponse a response.BusinessTypeResponse
func BusinessTypeDTOToResponse(businessTypeDTO domain.BusinessTypeResponse) response.BusinessTypeResponse {
	return response.BusinessTypeResponse{
		ID:          businessTypeDTO.ID,
		Name:        businessTypeDTO.Name,
		Code:        businessTypeDTO.Code,
		Description: businessTypeDTO.Description,
		Icon:        businessTypeDTO.Icon,
		IsActive:    businessTypeDTO.IsActive,
		CreatedAt:   businessTypeDTO.CreatedAt,
		UpdatedAt:   businessTypeDTO.UpdatedAt,
	}
}

// BusinessTypeDTOsToResponse convierte un slice de DTOs BusinessTypeResponse a slice de response.BusinessTypeResponse
func BusinessTypeDTOsToResponse(businessTypeDTOs []domain.BusinessTypeResponse) []response.BusinessTypeResponse {
	responses := make([]response.BusinessTypeResponse, len(businessTypeDTOs))
	for i, businessTypeDTO := range businessTypeDTOs {
		responses[i] = BusinessTypeDTOToResponse(businessTypeDTO)
	}
	return responses
}

// BuildCreateBusinessTypeResponseFromDTO construye la respuesta completa para crear un tipo de negocio desde un DTO
func BuildCreateBusinessTypeResponseFromDTO(businessTypeDTO *domain.BusinessTypeResponse, message string) response.CreateBusinessTypeResponse {
	return response.CreateBusinessTypeResponse{
		Success: true,
		Message: message,
		Data:    BusinessTypeDTOToResponse(*businessTypeDTO),
	}
}

// BusinessTypeToResponse convierte una entidad BusinessType del dominio a BusinessTypeResponse
func BusinessTypeToResponse(businessType domain.BusinessType) response.BusinessTypeResponse {
	return response.BusinessTypeResponse{
		ID:          businessType.ID,
		Name:        businessType.Name,
		Code:        businessType.Code,
		Description: businessType.Description,
		Icon:        businessType.Icon,
		IsActive:    businessType.DeletedAt == nil,
		CreatedAt:   businessType.CreatedAt,
		UpdatedAt:   businessType.UpdatedAt,
	}
}

// BusinessTypesToResponse convierte un slice de entidades BusinessType a slice de BusinessTypeResponse
func BusinessTypesToResponse(businessTypes []domain.BusinessType) []response.BusinessTypeResponse {
	responses := make([]response.BusinessTypeResponse, len(businessTypes))
	for i, businessType := range businessTypes {
		responses[i] = BusinessTypeToResponse(businessType)
	}
	return responses
}

// BuildGetBusinessTypesResponse construye la respuesta completa para obtener múltiples tipos de negocio
func BuildGetBusinessTypesResponse(businessTypes []domain.BusinessType, message string) response.GetBusinessTypesResponse {
	return response.GetBusinessTypesResponse{
		Success: true,
		Message: message,
		Data:    BusinessTypesToResponse(businessTypes),
	}
}

// BuildGetBusinessTypeResponse construye la respuesta completa para obtener un tipo de negocio
func BuildGetBusinessTypeResponse(businessType domain.BusinessType, message string) response.GetBusinessTypeResponse {
	return response.GetBusinessTypeResponse{
		Success: true,
		Message: message,
		Data:    BusinessTypeToResponse(businessType),
	}
}

// BuildUpdateBusinessTypeResponse construye la respuesta completa para actualizar un tipo de negocio
func BuildUpdateBusinessTypeResponse(businessType domain.BusinessType, message string) response.UpdateBusinessTypeResponse {
	return response.UpdateBusinessTypeResponse{
		Success: true,
		Message: message,
		Data:    BusinessTypeToResponse(businessType),
	}
}

// BuildDeleteBusinessTypeResponse construye la respuesta completa para eliminar un tipo de negocio
func BuildDeleteBusinessTypeResponse(message string) response.DeleteBusinessTypeResponse {
	return response.DeleteBusinessTypeResponse{
		Success: true,
		Message: message,
	}
}

// BuildErrorResponse construye una respuesta de error
func BuildErrorResponse(error string, message string) response.ErrorResponse {
	return response.ErrorResponse{
		Success: false,
		Error:   error,
		Message: message,
	}
}

// BuildGetBusinessTypesResponseFromDTOs construye la respuesta completa para obtener múltiples tipos de negocio desde DTOs
func BuildGetBusinessTypesResponseFromDTOs(businessTypeDTOs []domain.BusinessTypeResponse, message string) response.GetBusinessTypesResponse {
	return response.GetBusinessTypesResponse{
		Success: true,
		Message: message,
		Data:    BusinessTypeDTOsToResponse(businessTypeDTOs),
	}
}

// BuildGetBusinessTypeResponseFromDTO construye la respuesta completa para obtener un tipo de negocio desde un DTO
func BuildGetBusinessTypeResponseFromDTO(businessTypeDTO *domain.BusinessTypeResponse, message string) response.GetBusinessTypeResponse {
	return response.GetBusinessTypeResponse{
		Success: true,
		Message: message,
		Data:    BusinessTypeDTOToResponse(*businessTypeDTO),
	}
}

// BuildUpdateBusinessTypeResponseFromDTO construye la respuesta completa para actualizar un tipo de negocio desde un DTO
func BuildUpdateBusinessTypeResponseFromDTO(businessTypeDTO *domain.BusinessTypeResponse, message string) response.UpdateBusinessTypeResponse {
	return response.UpdateBusinessTypeResponse{
		Success: true,
		Message: message,
		Data:    BusinessTypeDTOToResponse(*businessTypeDTO),
	}
}
