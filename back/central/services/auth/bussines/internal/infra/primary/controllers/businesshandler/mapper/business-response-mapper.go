package mapper

import (
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/request"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/response"
)

// RequestToDTO convierte request.BusinessRequest a dtos.BusinessRequest
func RequestToDTO(req request.BusinessRequest) domain.BusinessRequest {
	return domain.BusinessRequest{
		Name:               req.Name,
		Code:               req.Code,
		BusinessTypeID:     req.BusinessTypeID,
		Timezone:           req.Timezone,
		Address:            req.Address,
		Description:        req.Description,
		LogoFile:           req.LogoFile,
		PrimaryColor:       req.PrimaryColor,
		SecondaryColor:     req.SecondaryColor,
		TertiaryColor:      req.TertiaryColor,
		QuaternaryColor:    req.QuaternaryColor,
		NavbarImageFile:    req.NavbarImageFile,
		CustomDomain:       req.CustomDomain,
		IsActive:           req.IsActive,
		EnableDelivery:     req.EnableDelivery,
		EnablePickup:       req.EnablePickup,
		EnableReservations: req.EnableReservations,
	}
}

// RequestToUpdateDTO convierte request.BusinessRequest a dtos.UpdateBusinessRequest
func RequestToUpdateDTO(req request.BusinessRequest) domain.UpdateBusinessRequest {
	return domain.UpdateBusinessRequest{
		Name:               &req.Name,
		Code:               &req.Code,
		BusinessTypeID:     &req.BusinessTypeID,
		Timezone:           &req.Timezone,
		Address:            &req.Address,
		Description:        &req.Description,
		LogoFile:           req.LogoFile,
		PrimaryColor:       &req.PrimaryColor,
		SecondaryColor:     &req.SecondaryColor,
		TertiaryColor:      &req.TertiaryColor,
		QuaternaryColor:    &req.QuaternaryColor,
		NavbarImageFile:    req.NavbarImageFile,
		CustomDomain:       &req.CustomDomain,
		IsActive:           &req.IsActive,
		EnableDelivery:     &req.EnableDelivery,
		EnablePickup:       &req.EnablePickup,
		EnableReservations: &req.EnableReservations,
	}
}

// UpdateRequestToUpdateDTO convierte request.UpdateBusinessRequest a dtos.UpdateBusinessRequest
func UpdateRequestToUpdateDTO(req request.UpdateBusinessRequest) domain.UpdateBusinessRequest {
	return domain.UpdateBusinessRequest{
		Name:               &req.Name,
		Code:               &req.Code,
		BusinessTypeID:     &req.BusinessTypeID,
		Timezone:           &req.Timezone,
		Address:            &req.Address,
		Description:        &req.Description,
		LogoFile:           req.LogoFile,
		PrimaryColor:       &req.PrimaryColor,
		SecondaryColor:     &req.SecondaryColor,
		TertiaryColor:      &req.TertiaryColor,
		QuaternaryColor:    &req.QuaternaryColor,
		NavbarImageFile:    req.NavbarImageFile,
		CustomDomain:       &req.CustomDomain,
		IsActive:           &req.IsActive,
		EnableDelivery:     &req.EnableDelivery,
		EnablePickup:       &req.EnablePickup,
		EnableReservations: &req.EnableReservations,
	}
}

// BusinessToResponse convierte una entidad Business del dominio a BusinessResponse
func BusinessToResponse(business domain.Business) response.BusinessResponse {
	return response.BusinessResponse{
		ID:              business.ID,
		Name:            business.Name,
		Description:     business.Description,
		Address:         business.Address,
		Phone:           "", // Campo no disponible en la entidad
		Email:           "", // Campo no disponible en la entidad
		Website:         "", // Campo no disponible en la entidad
		LogoURL:         business.LogoURL,
		PrimaryColor:    business.PrimaryColor,
		SecondaryColor:  business.SecondaryColor,
		TertiaryColor:   business.TertiaryColor,
		QuaternaryColor: business.QuaternaryColor,
		NavbarImageURL:  business.NavbarImageURL,
		IsActive:        business.DeletedAt == nil,
		BusinessTypeID:  business.BusinessTypeID,
		BusinessType:    "", // Se llenará desde el handler
		CreatedAt:       business.CreatedAt,
		UpdatedAt:       business.UpdatedAt,
	}
}

// BusinessDTOToResponse convierte un DTO BusinessResponse a response.BusinessResponse
func BusinessDTOToResponse(businessDTO domain.BusinessResponse) response.BusinessResponse {
	return response.BusinessResponse{
		ID:              businessDTO.ID,
		Name:            businessDTO.Name,
		Description:     businessDTO.Description,
		Address:         businessDTO.Address,
		Phone:           "", // Campo no disponible en el DTO
		Email:           "", // Campo no disponible en el DTO
		Website:         "", // Campo no disponible en el DTO
		LogoURL:         businessDTO.LogoURL,
		PrimaryColor:    businessDTO.PrimaryColor,
		SecondaryColor:  businessDTO.SecondaryColor,
		TertiaryColor:   businessDTO.TertiaryColor,
		QuaternaryColor: businessDTO.QuaternaryColor,
		NavbarImageURL:  businessDTO.NavbarImageURL,
		IsActive:        businessDTO.IsActive,
		BusinessTypeID:  businessDTO.BusinessType.ID,
		BusinessType:    businessDTO.BusinessType.Name,
		CreatedAt:       businessDTO.CreatedAt,
		UpdatedAt:       businessDTO.UpdatedAt,
	}
}

// BusinessDTOToDetailResponse convierte un DTO BusinessResponse a response.BusinessDetailResponse
func BusinessDTOToDetailResponse(businessDTO domain.BusinessResponse) response.BusinessDetailResponse {
	return response.BusinessDetailResponse{
		ID:   businessDTO.ID,
		Name: businessDTO.Name,
		Code: businessDTO.Code,
		BusinessType: response.BusinessTypeDetailResponse{
			ID:          businessDTO.BusinessType.ID,
			Name:        businessDTO.BusinessType.Name,
			Code:        businessDTO.BusinessType.Code,
			Description: businessDTO.BusinessType.Description,
			Icon:        businessDTO.BusinessType.Icon,
			IsActive:    businessDTO.BusinessType.IsActive,
			CreatedAt:   businessDTO.BusinessType.CreatedAt,
			UpdatedAt:   businessDTO.BusinessType.UpdatedAt,
		},
		Timezone:           businessDTO.Timezone,
		Address:            businessDTO.Address,
		Description:        businessDTO.Description,
		LogoURL:            businessDTO.LogoURL,
		PrimaryColor:       businessDTO.PrimaryColor,
		SecondaryColor:     businessDTO.SecondaryColor,
		TertiaryColor:      businessDTO.TertiaryColor,
		QuaternaryColor:    businessDTO.QuaternaryColor,
		NavbarImageURL:     businessDTO.NavbarImageURL,
		CustomDomain:       businessDTO.CustomDomain,
		IsActive:           businessDTO.IsActive,
		EnableDelivery:     businessDTO.EnableDelivery,
		EnablePickup:       businessDTO.EnablePickup,
		EnableReservations: businessDTO.EnableReservations,
		CreatedAt:          businessDTO.CreatedAt,
		UpdatedAt:          businessDTO.UpdatedAt,
	}
}

// BusinessesToResponse convierte un slice de entidades Business a slice de BusinessResponse
func BusinessesToResponse(businesses []domain.Business) []response.BusinessResponse {
	responses := make([]response.BusinessResponse, len(businesses))
	for i, business := range businesses {
		responses[i] = BusinessToResponse(business)
	}
	return responses
}

// BusinessDTOsToResponse convierte un slice de DTOs BusinessResponse a slice de response.BusinessResponse
func BusinessDTOsToResponse(businessDTOs []domain.BusinessResponse) []response.BusinessResponse {
	responses := make([]response.BusinessResponse, len(businessDTOs))
	for i, businessDTO := range businessDTOs {
		responses[i] = BusinessDTOToResponse(businessDTO)
	}
	return responses
}

// BuildGetBusinessesResponse construye la respuesta completa para obtener múltiples negocios
func BuildGetBusinessesResponse(businesses []domain.Business, message string) response.GetBusinessesResponse {
	return response.GetBusinessesResponse{
		Success: true,
		Message: message,
		Data:    BusinessesToResponse(businesses),
	}
}

// BuildGetBusinessesResponseFromDTOs construye la respuesta completa para obtener múltiples negocios desde DTOs
func BuildGetBusinessesResponseFromDTOs(businessDTOs []domain.BusinessResponse, message string) response.GetBusinessesResponse {
	return response.GetBusinessesResponse{
		Success: true,
		Message: message,
		Data:    BusinessDTOsToResponse(businessDTOs),
	}
}

// BuildGetBusinessesResponseWithPagination construye la respuesta completa para obtener múltiples negocios con paginación
func BuildGetBusinessesResponseWithPagination(businessDTOs []domain.BusinessResponse, message string, page, limit int, total int64) response.GetBusinessesResponse {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasNext := page < totalPages
	hasPrev := page > 1

	return response.GetBusinessesResponse{
		Success: true,
		Message: message,
		Data:    BusinessDTOsToResponse(businessDTOs),
		Pagination: &response.PaginationInfo{
			CurrentPage: page,
			PerPage:     limit,
			Total:       total,
			LastPage:    totalPages,
			HasNext:     hasNext,
			HasPrev:     hasPrev,
		},
	}
}

// BuildGetBusinessResponse construye la respuesta completa para obtener un negocio
func BuildGetBusinessResponse(business domain.Business, message string) response.GetBusinessResponse {
	return response.GetBusinessResponse{
		Success: true,
		Message: message,
		Data:    BusinessToResponse(business),
	}
}

// BuildCreateBusinessResponse construye la respuesta completa para crear un negocio
func BuildCreateBusinessResponse(business domain.Business, message string) response.CreateBusinessResponse {
	return response.CreateBusinessResponse{
		Success: true,
		Message: message,
		Data:    BusinessToResponse(business),
	}
}

// BuildUpdateBusinessResponse construye la respuesta completa para actualizar un negocio
func BuildUpdateBusinessResponse(business domain.Business, message string) response.UpdateBusinessResponse {
	return response.UpdateBusinessResponse{
		Success: true,
		Message: message,
		Data:    BusinessToResponse(business),
	}
}

// BuildDeleteBusinessResponse construye la respuesta completa para eliminar un negocio
func BuildDeleteBusinessResponse(message string) response.DeleteBusinessResponse {
	return response.DeleteBusinessResponse{
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

// BuildCreateBusinessResponseFromDTO construye la respuesta completa para crear un negocio desde un DTO
func BuildCreateBusinessResponseFromDTO(businessDTO *domain.BusinessResponse, message string) response.CreateBusinessResponse {
	return response.CreateBusinessResponse{
		Success: true,
		Message: message,
		Data:    BusinessDTOToResponse(*businessDTO),
	}
}

// BuildGetBusinessResponseFromDTO construye la respuesta completa para obtener un negocio desde un DTO
func BuildGetBusinessResponseFromDTO(businessDTO *domain.BusinessResponse, message string) response.GetBusinessResponse {
	return response.GetBusinessResponse{
		Success: true,
		Message: message,
		Data:    BusinessDTOToResponse(*businessDTO),
	}
}

// BuildGetBusinessByIDResponseFromDTO construye la respuesta completa para obtener un negocio con toda su información desde un DTO
func BuildGetBusinessByIDResponseFromDTO(businessDTO *domain.BusinessResponse, message string) response.GetBusinessByIDResponse {
	return response.GetBusinessByIDResponse{
		Success: true,
		Message: message,
		Data:    BusinessDTOToDetailResponse(*businessDTO),
	}
}

// BuildUpdateBusinessResponseFromDTO construye la respuesta completa para actualizar un negocio desde un DTO
func BuildUpdateBusinessResponseFromDTO(businessDTO *domain.BusinessResponse, message string) response.UpdateBusinessResponse {
	return response.UpdateBusinessResponse{
		Success: true,
		Message: message,
		Data:    BusinessDTOToResponse(*businessDTO),
	}
}
