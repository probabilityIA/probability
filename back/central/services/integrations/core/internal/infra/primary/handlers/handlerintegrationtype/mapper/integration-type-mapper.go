package mapper

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
	"gorm.io/datatypes"
)

// ToIntegrationTypeResponse convierte domain.IntegrationType a IntegrationTypeResponse
func ToIntegrationTypeResponse(integrationType *domain.IntegrationType) response.IntegrationTypeResponse {

	return response.IntegrationTypeResponse{
		ID:                integrationType.ID,
		Name:              integrationType.Name,
		Code:              integrationType.Code,
		Description:       integrationType.Description,
		Icon:              integrationType.Icon,
		Category:          integrationType.Category,
		IsActive:          integrationType.IsActive,
		ConfigSchema:      integrationType.ConfigSchema,
		CredentialsSchema: integrationType.CredentialsSchema,
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
		Category:          req.Category,
		IsActive:          req.IsActive,
		ConfigSchema:      configSchema,
		CredentialsSchema: credentialsSchema,
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
	if req.Category != nil {
		dto.Category = req.Category
	}
	if req.IsActive != nil {
		dto.IsActive = req.IsActive
	}
	if req.ConfigSchema != nil {
		configBytes, _ := json.Marshal(*req.ConfigSchema)
		configJSON := datatypes.JSON(configBytes)
		dto.ConfigSchema = &configJSON
	}

	return dto
}
