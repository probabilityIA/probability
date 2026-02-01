package mappers

import (
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/response"
)

// ProviderToResponse convierte entidad de dominio a response
func ProviderToResponse(provider *entities.Provider, providerTypeCode string) *response.Provider {
	return &response.Provider{
		ID:               provider.ID,
		CreatedAt:        provider.CreatedAt,
		UpdatedAt:        provider.UpdatedAt,
		Name:             provider.Name,
		Description:      provider.Description,
		ProviderTypeCode: providerTypeCode,
		BusinessID:       provider.BusinessID,
		Config:           provider.Config,
		IsActive:         provider.IsActive,
		IsDefault:        provider.IsDefault,
	}
}

// ProvidersToResponse convierte lista de proveedores a response de listado
func ProvidersToResponse(providers []*entities.Provider, providerTypes map[uint]string, total int64, page, pageSize int) *response.ProviderList {
	items := make([]response.Provider, 0, len(providers))

	for _, provider := range providers {
		typeCode := providerTypes[provider.ProviderTypeID]
		if typeCode == "" {
			typeCode = "unknown"
		}
		items = append(items, *ProviderToResponse(provider, typeCode))
	}

	return &response.ProviderList{
		Items:      items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}
}

// ProviderTypeToResponse convierte tipo de proveedor a response
func ProviderTypeToResponse(providerType *entities.ProviderType) *response.ProviderType {
	return &response.ProviderType{
		ID:                 providerType.ID,
		Code:               providerType.Code,
		Name:               providerType.Name,
		Description:        providerType.Description,
		Icon:               providerType.Icon,
		ImageURL:           providerType.ImageURL,
		ApiBaseURL:         providerType.ApiBaseURL,
		DocumentationURL:   providerType.DocumentationURL,
		SupportedCountries: providerType.SupportedCountries,
	}
}

// ProviderTypesToResponse convierte lista de tipos a response
func ProviderTypesToResponse(providerTypes []*entities.ProviderType) []response.ProviderType {
	items := make([]response.ProviderType, 0, len(providerTypes))

	for _, pt := range providerTypes {
		items = append(items, *ProviderTypeToResponse(pt))
	}

	return items
}
