package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/errors"
)

// UpdateProvider actualiza un proveedor existente
func (uc *useCase) UpdateProvider(ctx context.Context, id uint, dto *dtos.UpdateProviderDTO) (*entities.Provider, error) {
	uc.log.Info(ctx).Uint("provider_id", id).Msg("Updating Softpymes provider")

	// 1. Obtener proveedor existente
	provider, err := uc.providerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrProviderNotFound
	}

	// 2. Actualizar solo los campos proporcionados
	if dto.Name != nil {
		provider.Name = *dto.Name
	}

	if dto.Description != nil {
		provider.Description = *dto.Description
	}

	if dto.Config != nil {
		provider.Config = dto.Config
	}

	if dto.IsActive != nil {
		provider.IsActive = *dto.IsActive
	}

	if dto.IsDefault != nil {
		// Si se marca como default, verificar que no haya otro
		if *dto.IsDefault {
			defaultProvider, err := uc.providerRepo.GetDefaultByBusiness(ctx, provider.BusinessID)
			if err == nil && defaultProvider != nil && defaultProvider.ID != id {
				return nil, errors.ErrDefaultProviderExists
			}
		}
		provider.IsDefault = *dto.IsDefault
	}

	// 3. Actualizar credenciales si fueron proporcionadas
	if dto.Credentials != nil {
		provider.Credentials = dto.Credentials
	}

	// 4. Guardar cambios
	if err := uc.providerRepo.Update(ctx, provider); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update provider")
		return nil, err
	}

	uc.log.Info(ctx).Uint("provider_id", provider.ID).Msg("Provider updated successfully")
	return provider, nil
}
