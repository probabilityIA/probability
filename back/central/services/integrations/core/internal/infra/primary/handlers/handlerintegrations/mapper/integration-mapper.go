package mapper

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
	"gorm.io/datatypes"
)

// ToCreateIntegrationDTO convierte CreateIntegrationRequest a CreateIntegrationDTO
func ToCreateIntegrationDTO(req request.CreateIntegrationRequest, createdByID uint) domain.CreateIntegrationDTO {
	var configJSON datatypes.JSON
	if req.Config != nil {
		configBytes, _ := json.Marshal(req.Config)
		configJSON = configBytes
	}

	return domain.CreateIntegrationDTO{
		Name:              req.Name,
		Code:              req.Code,
		IntegrationTypeID: req.IntegrationTypeID,
		Category:          req.Category,
		BusinessID:        req.BusinessID,
		IsActive:          req.IsActive,
		IsDefault:         req.IsDefault,
		Config:            configJSON,
		Credentials:       req.Credentials,
		Description:       req.Description,
		CreatedByID:       createdByID,
	}
}

// ToUpdateIntegrationDTO convierte UpdateIntegrationRequest a UpdateIntegrationDTO
func ToUpdateIntegrationDTO(req request.UpdateIntegrationRequest, updatedByID uint) domain.UpdateIntegrationDTO {
	dto := domain.UpdateIntegrationDTO{
		UpdatedByID: updatedByID,
	}

	if req.Name != nil {
		dto.Name = req.Name
	}
	if req.Code != nil {
		dto.Code = req.Code
	}
	if req.IntegrationTypeID != nil {
		dto.IntegrationTypeID = req.IntegrationTypeID
	}
	if req.IsActive != nil {
		dto.IsActive = req.IsActive
	}
	if req.IsDefault != nil {
		dto.IsDefault = req.IsDefault
	}
	if req.Description != nil {
		dto.Description = req.Description
	}
	if req.Config != nil {
		configBytes, _ := json.Marshal(*req.Config)
		configJSON := datatypes.JSON(configBytes)
		dto.Config = &configJSON
	}
	if req.Credentials != nil {
		dto.Credentials = req.Credentials
	}

	return dto
}

// ToIntegrationResponse convierte domain.Integration a IntegrationResponse
func ToIntegrationResponse(integration *domain.Integration) response.IntegrationResponse {
	var config map[string]interface{}

	resp := response.IntegrationResponse{
		ID:                integration.ID,
		Name:              integration.Name,
		Code:              integration.Code,
		IntegrationTypeID: integration.IntegrationTypeID,
		Category:          integration.Category,
		BusinessID:        integration.BusinessID,
		IsActive:          integration.IsActive,
		IsDefault:         integration.IsDefault,
		Config:            config,
		Description:       integration.Description,
		CreatedByID:       integration.CreatedByID,
		UpdatedByID:       integration.UpdatedByID,
		CreatedAt:         integration.CreatedAt,
		UpdatedAt:         integration.UpdatedAt,
	}

	// Incluir información del tipo de integración si está cargado
	if integration.IntegrationType != nil {
		resp.IntegrationType = &response.IntegrationTypeInfo{
			ID:   integration.IntegrationType.ID,
			Name: integration.IntegrationType.Name,
			Code: integration.IntegrationType.Code,
		}
	}

	return resp
}

// ToIntegrationListResponse convierte lista de integraciones a IntegrationListResponse
func ToIntegrationListResponse(integrations []*domain.Integration, total int64, page, pageSize int) response.IntegrationListResponse {
	responses := make([]response.IntegrationResponse, len(integrations))
	for i, integration := range integrations {
		responses[i] = ToIntegrationResponse(integration)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return response.IntegrationListResponse{
		Success:    true,
		Message:    "Integraciones obtenidas exitosamente",
		Data:       responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ToIntegrationFilters convierte GetIntegrationsRequest a IntegrationFilters
func ToIntegrationFilters(req request.GetIntegrationsRequest) domain.IntegrationFilters {
	return domain.IntegrationFilters{
		Page:                req.Page,
		PageSize:            req.PageSize,
		IntegrationTypeID:   req.IntegrationTypeID,
		IntegrationTypeCode: req.IntegrationTypeCode,
		Category:            req.Category,
		BusinessID:          req.BusinessID,
		IsActive:            req.IsActive,
		Search:              req.Search,
	}
}
