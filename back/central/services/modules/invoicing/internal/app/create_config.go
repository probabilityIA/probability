package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// CreateConfig crea una nueva configuración de facturación
func (uc *useCase) CreateConfig(ctx context.Context, dto *dtos.CreateConfigDTO) (*entities.InvoicingConfig, error) {
	uc.log.Info(ctx).Uint("integration_id", dto.IntegrationID).Msg("Creating invoicing config")

	// 1. Verificar si ya existe config para esta integración
	existingConfig, err := uc.repo.GetConfigByIntegration(ctx, dto.IntegrationID)
	if err == nil && existingConfig != nil {
		if existingConfig.Enabled {
			// Config activa: bloquear — el usuario debe desactivarla primero
			uc.log.Warn(ctx).
				Uint("integration_id", dto.IntegrationID).
				Uint("existing_config_id", existingConfig.ID).
				Msg("Active config already exists for this integration")
			return nil, errors.ErrConfigAlreadyExists
		}
		// Config inactiva: actualizarla con los nuevos datos (upsert)
		uc.log.Info(ctx).
			Uint("integration_id", dto.IntegrationID).
			Uint("existing_config_id", existingConfig.ID).
			Msg("Updating existing disabled config with new values")

		invoicingIntegrationID := &dto.InvoicingIntegrationID
		existingConfig.InvoicingIntegrationID = invoicingIntegrationID
		existingConfig.Enabled = dto.Enabled
		existingConfig.AutoInvoice = dto.AutoInvoice
		existingConfig.Filters = dto.Filters
		existingConfig.InvoiceConfig = dto.InvoiceConfig
		if dto.Description != nil {
			existingConfig.Description = *dto.Description
		}

		// Si se activa, verificar que no hay otro config activo en el negocio
		if dto.Enabled {
			activeConfig, err := uc.repo.GetEnabledConfigByBusiness(ctx, dto.BusinessID)
			if err != nil {
				return nil, err
			}
			if activeConfig != nil && activeConfig.ID != existingConfig.ID {
				uc.log.Warn(ctx).
					Uint("business_id", dto.BusinessID).
					Uint("active_config_id", activeConfig.ID).
					Msg("Business already has an active invoicing config")
				return nil, errors.ErrActiveInvoicingConfigExists
			}
		}

		if err := uc.repo.UpdateInvoicingConfig(ctx, existingConfig); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Failed to update existing config")
			return nil, err
		}
		uc.log.Info(ctx).Uint("config_id", existingConfig.ID).Msg("Existing config updated successfully")
		return existingConfig, nil
	}

	// 2. Si se crea activa, verificar que no hay otro config activo para este negocio
	if dto.Enabled {
		activeConfig, err := uc.repo.GetEnabledConfigByBusiness(ctx, dto.BusinessID)
		if err != nil {
			return nil, err
		}
		if activeConfig != nil {
			uc.log.Warn(ctx).
				Uint("business_id", dto.BusinessID).
				Uint("active_config_id", activeConfig.ID).
				Msg("Business already has an active invoicing config")
			return nil, errors.ErrActiveInvoicingConfigExists
		}
	}

	// 2. TODO: Validar que el proveedor existe y está activo usando integrationCore
	// Por ahora omitimos la validación para que compile
	// provider, err := uc.integrationCore.GetIntegrationByID(ctx, fmt.Sprintf("%d", dto.InvoicingProviderID))
	// if err != nil {
	// 	return nil, errors.ErrProviderNotFound
	// }

	// 3. Crear entidad
	invoicingIntegrationID := &dto.InvoicingIntegrationID

	description := ""
	if dto.Description != nil {
		description = *dto.Description
	}

	config := &entities.InvoicingConfig{
		BusinessID:             dto.BusinessID,
		IntegrationID:          dto.IntegrationID,
		InvoicingIntegrationID: invoicingIntegrationID,
		Enabled:                dto.Enabled,
		AutoInvoice:            dto.AutoInvoice,
		Filters:                dto.Filters,
		InvoiceConfig:          dto.InvoiceConfig,
		Description:            description,
		CreatedByID:            dto.CreatedByUserID,
	}

	// 4. Guardar en BD
	if err := uc.repo.CreateInvoicingConfig(ctx, config); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create config")
		return nil, err
	}

	uc.log.Info(ctx).Uint("config_id", config.ID).Msg("Config created successfully")
	return config, nil
}

// UpdateConfig actualiza una configuración existente
func (uc *useCase) UpdateConfig(ctx context.Context, id uint, dto *dtos.UpdateConfigDTO) (*entities.InvoicingConfig, error) {
	uc.log.Info(ctx).Uint("config_id", id).Msg("Updating invoicing config")

	// Obtener config existente
	config, err := uc.repo.GetInvoicingConfigByID(ctx, id)
	if err != nil {
		return nil, errors.ErrConfigNotFound
	}

	// Si se está activando, verificar que no hay otro config activo para este negocio
	if dto.Enabled != nil && *dto.Enabled && !config.Enabled {
		activeConfig, err := uc.repo.GetEnabledConfigByBusiness(ctx, config.BusinessID)
		if err != nil {
			return nil, err
		}
		if activeConfig != nil && activeConfig.ID != id {
			uc.log.Warn(ctx).
				Uint("business_id", config.BusinessID).
				Uint("active_config_id", activeConfig.ID).
				Uint("requested_config_id", id).
				Msg("Cannot activate config: business already has an active invoicing config")
			return nil, errors.ErrActiveInvoicingConfigExists
		}
	}

	// Actualizar solo los campos proporcionados
	if dto.Enabled != nil {
		config.Enabled = *dto.Enabled
	}

	if dto.AutoInvoice != nil {
		config.AutoInvoice = *dto.AutoInvoice
	}

	if dto.InvoicingIntegrationID != nil {
		config.InvoicingIntegrationID = dto.InvoicingIntegrationID
	}

	if dto.Filters != nil {
		config.Filters = dto.Filters
	}

	// Guardar cambios
	if err := uc.repo.UpdateInvoicingConfig(ctx, config); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update config")
		return nil, err
	}

	uc.log.Info(ctx).Uint("config_id", config.ID).Msg("Config updated successfully")
	return config, nil
}

// GetConfig obtiene una configuración por ID
func (uc *useCase) GetConfig(ctx context.Context, id uint) (*entities.InvoicingConfig, error) {
	return uc.repo.GetInvoicingConfigByID(ctx, id)
}

// ListConfigs lista configuraciones de un negocio
func (uc *useCase) ListConfigs(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error) {
	return uc.repo.ListInvoicingConfigs(ctx, businessID)
}
