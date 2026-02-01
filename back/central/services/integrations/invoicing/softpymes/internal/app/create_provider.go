package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/errors"
)

// CreateProvider crea un nuevo proveedor de facturación Softpymes
func (uc *useCase) CreateProvider(ctx context.Context, dto *dtos.CreateProviderDTO) (*entities.Provider, error) {
	uc.log.Info(ctx).Uint("business_id", dto.BusinessID).Str("type", dto.ProviderTypeCode).Msg("Creating Softpymes provider")

	// 1. Validar que el tipo de proveedor existe
	providerType, err := uc.providerTypeRepo.GetByCode(ctx, dto.ProviderTypeCode)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("code", dto.ProviderTypeCode).Msg("Provider type not found")
		return nil, errors.ErrProviderTypeNotFound
	}

	if !providerType.IsActive {
		return nil, errors.ErrProviderTypeInactive
	}

	// 2. Validar que no existe otro proveedor del mismo tipo para este negocio
	existing, err := uc.providerRepo.GetByBusinessAndType(ctx, dto.BusinessID, dto.ProviderTypeCode)
	if err == nil && existing != nil {
		return nil, errors.ErrProviderAlreadyExists
	}

	// 3. Si es proveedor por defecto, verificar que no haya otro
	if dto.IsDefault {
		defaultProvider, err := uc.providerRepo.GetDefaultByBusiness(ctx, dto.BusinessID)
		if err == nil && defaultProvider != nil {
			return nil, errors.ErrDefaultProviderExists
		}
	}

	// 4. Crear entidad
	description := ""
	if dto.Description != nil {
		description = *dto.Description
	}

	provider := &entities.Provider{
		BusinessID:     dto.BusinessID,
		ProviderTypeID: providerType.ID,
		Name:           dto.Name,
		Description:    description,
		IsActive:       true,
		IsDefault:      dto.IsDefault,
		Config:         dto.Config,
		Credentials:    dto.Credentials, // Credenciales sin encriptar (se encriptan en integrationCore)
		CreatedByID:    dto.CreatedByUserID,
	}

	// 5. Guardar en BD
	if err := uc.providerRepo.Create(ctx, provider); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create provider")
		return nil, err
	}

	uc.log.Info(ctx).Uint("provider_id", provider.ID).Msg("Provider created successfully")
	return provider, nil
}
