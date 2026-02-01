package mappers

import (
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/request"
)

// CreateProviderRequestToDTO convierte request a DTO de dominio
func CreateProviderRequestToDTO(req *request.CreateProvider) *dtos.CreateProviderDTO {
	return &dtos.CreateProviderDTO{
		Name:             req.Name,
		ProviderTypeCode: req.ProviderTypeCode,
		BusinessID:       req.BusinessID,
		Description:      req.Description,
		Config:           req.Config,
		Credentials:      req.Credentials,
		IsDefault:        req.IsDefault,
		CreatedByUserID:  req.CreatedByUserID,
	}
}

// UpdateProviderRequestToDTO convierte request a DTO de dominio
func UpdateProviderRequestToDTO(req *request.UpdateProvider) *dtos.UpdateProviderDTO {
	dto := &dtos.UpdateProviderDTO{}

	if req.Name != nil {
		dto.Name = req.Name
	}

	if req.Description != nil {
		dto.Description = req.Description
	}

	if req.Config != nil {
		dto.Config = *req.Config
	}

	if req.Credentials != nil {
		dto.Credentials = *req.Credentials
	}

	if req.IsActive != nil {
		dto.IsActive = req.IsActive
	}

	if req.IsDefault != nil {
		dto.IsDefault = req.IsDefault
	}

	return dto
}
