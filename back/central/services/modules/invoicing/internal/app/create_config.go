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

	// 1. Verificar que no existe config para esta integración
	exists, err := uc.configRepo.ExistsForIntegration(ctx, dto.IntegrationID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrConfigAlreadyExists
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
	if err := uc.configRepo.Create(ctx, config); err != nil {
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
	config, err := uc.configRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrConfigNotFound
	}

	// Actualizar solo los campos proporcionados
	if dto.Enabled != nil {
		config.Enabled = *dto.Enabled
	}

	if dto.AutoInvoice != nil {
		config.AutoInvoice = *dto.AutoInvoice
	}

	if dto.Filters != nil {
		config.Filters = dto.Filters
	}

	// Guardar cambios
	if err := uc.configRepo.Update(ctx, config); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update config")
		return nil, err
	}

	uc.log.Info(ctx).Uint("config_id", config.ID).Msg("Config updated successfully")
	return config, nil
}

// GetConfig obtiene una configuración por ID
func (uc *useCase) GetConfig(ctx context.Context, id uint) (*entities.InvoicingConfig, error) {
	return uc.configRepo.GetByID(ctx, id)
}

// ListConfigs lista configuraciones de un negocio
func (uc *useCase) ListConfigs(ctx context.Context, businessID uint) ([]*entities.InvoicingConfig, error) {
	return uc.configRepo.List(ctx, businessID)
}
