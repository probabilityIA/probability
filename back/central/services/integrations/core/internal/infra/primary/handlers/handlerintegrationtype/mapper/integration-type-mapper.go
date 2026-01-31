package mapper

import (
	"encoding/json"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
	"gorm.io/datatypes"
)

// ToIntegrationTypeResponse convierte domain.IntegrationType a IntegrationTypeResponse
// imageURLBase es la URL base de S3 para construir URLs completas
func ToIntegrationTypeResponse(integrationType *domain.IntegrationType, imageURLBase string) response.IntegrationTypeResponse {
	imageURL := ""
	if integrationType.ImageURL != "" {
		// Si ya es una URL completa, usarla directamente
		if imageURLBase != "" && !strings.HasPrefix(integrationType.ImageURL, "http") {
			imageURL = strings.TrimRight(imageURLBase, "/") + "/" + strings.TrimLeft(integrationType.ImageURL, "/")
		} else {
			imageURL = integrationType.ImageURL
		}
	}

	var category *response.IntegrationCategoryResponse
	if integrationType.Category != nil {
		category = &response.IntegrationCategoryResponse{
			ID:          integrationType.Category.ID,
			Code:        integrationType.Category.Code,
			Name:        integrationType.Category.Name,
			Description: integrationType.Category.Description,
			Icon:        integrationType.Category.Icon,
			Color:       integrationType.Category.Color,
		}
	}

	return response.IntegrationTypeResponse{
		ID:                integrationType.ID,
		Name:              integrationType.Name,
		Code:              integrationType.Code,
		Description:       integrationType.Description,
		Icon:              integrationType.Icon,
		ImageURL:          imageURL,
		Category:          category,
		IsActive:          integrationType.IsActive,
		ConfigSchema:      integrationType.ConfigSchema,
		CredentialsSchema: integrationType.CredentialsSchema,
		SetupInstructions: integrationType.SetupInstructions,
		CreatedAt:         integrationType.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         integrationType.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToCreateIntegrationTypeDTO convierte CreateIntegrationTypeRequest a CreateIntegrationTypeDTO
func ToCreateIntegrationTypeDTO(req request.CreateIntegrationTypeRequest) domain.CreateIntegrationTypeDTO {
	var configSchema, credentialsSchema datatypes.JSON
	if req.ConfigSchema != nil {
		configBytes, _ := json.Marshal(req.ConfigSchema)
		configSchema = configBytes
	}

	return domain.CreateIntegrationTypeDTO{
		Name:              req.Name,
		Code:              req.Code,
		Description:       req.Description,
		Icon:              req.Icon,
		CategoryID:        req.CategoryID,
		IsActive:          req.IsActive,
		ConfigSchema:      configSchema,
		CredentialsSchema: credentialsSchema,
		ImageFile:         req.ImageFile,
	}
}

// ToUpdateIntegrationTypeDTO convierte UpdateIntegrationTypeRequest a UpdateIntegrationTypeDTO
func ToUpdateIntegrationTypeDTO(req request.UpdateIntegrationTypeRequest) domain.UpdateIntegrationTypeDTO {
	dto := domain.UpdateIntegrationTypeDTO{}

	if req.Name != nil {
		dto.Name = req.Name
	}
	if req.Code != nil {
		dto.Code = req.Code
	}
	if req.Description != nil {
		dto.Description = req.Description
	}
	if req.Icon != nil {
		dto.Icon = req.Icon
	}
	if req.CategoryID != nil {
		dto.CategoryID = req.CategoryID
	}
	if req.IsActive != nil {
		dto.IsActive = req.IsActive
	}
	if req.ConfigSchema != nil {
		configBytes, _ := json.Marshal(*req.ConfigSchema)
		configJSON := datatypes.JSON(configBytes)
		dto.ConfigSchema = &configJSON
	}
	if req.ImageFile != nil {
		dto.ImageFile = req.ImageFile
	}
	if req.RemoveImage != nil {
		dto.RemoveImage = *req.RemoveImage
	}

	return dto
}
