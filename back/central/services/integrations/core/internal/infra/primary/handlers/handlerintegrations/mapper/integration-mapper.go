package mapper

import (
	"encoding/json"
	"strings"

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
// imageURLBase es la URL base de S3 para construir URLs completas
func ToIntegrationResponse(integration *domain.Integration, imageURLBase string) response.IntegrationResponse {
	var config map[string]interface{}

	// Parsear Config desde datatypes.JSON ([]byte) a map[string]interface{}
	if len(integration.Config) > 0 {
		if err := json.Unmarshal(integration.Config, &config); err != nil {
			// Si falla el parseo, dejar config vacío
			config = make(map[string]interface{})
		}
	}

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
		imageURL := ""
		if integration.IntegrationType.ImageURL != "" {
			// Construir URL completa si es path relativo
			if imageURLBase != "" && !strings.HasPrefix(integration.IntegrationType.ImageURL, "http") {
				imageURL = strings.TrimRight(imageURLBase, "/") + "/" + strings.TrimLeft(integration.IntegrationType.ImageURL, "/")
			} else {
				imageURL = integration.IntegrationType.ImageURL
			}
		}
		resp.IntegrationType = &response.IntegrationTypeInfo{
			ID:       integration.IntegrationType.ID,
			Name:     integration.IntegrationType.Name,
			Code:     integration.IntegrationType.Code,
			ImageURL: imageURL,
		}
	}

	return resp
}

// ToIntegrationListResponse convierte lista de integraciones a IntegrationListResponse
func ToIntegrationListResponse(integrations []*domain.Integration, total int64, page, pageSize int, imageURLBase string) response.IntegrationListResponse {
	responses := make([]response.IntegrationResponse, len(integrations))
	for i, integration := range integrations {
		responses[i] = ToIntegrationResponse(integration, imageURLBase)
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
