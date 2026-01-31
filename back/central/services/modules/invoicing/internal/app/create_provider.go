package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// CreateProvider crea un nuevo proveedor de facturación
func (uc *useCase) CreateProvider(ctx context.Context, dto *dtos.CreateProviderDTO) (*entities.InvoicingProvider, error) {
	uc.log.Info(ctx).Uint("business_id", dto.BusinessID).Str("type", dto.ProviderTypeCode).Msg("Creating invoicing provider")

	// 1. Validar que el tipo de proveedor existe
	providerType, err := uc.providerTypeRepo.GetByCode(ctx, dto.ProviderTypeCode)
	if err != nil {
		return nil, errors.ErrProviderTypeNotFound
	}

	// 2. Encriptar credenciales
	encryptedCreds, err := uc.encryption.Encrypt(dto.Credentials)
	if err != nil {
		return nil, errors.ErrEncryptionFailed
	}

	// 3. Crear entidad
	provider := &entities.InvoicingProvider{
		BusinessID:     dto.BusinessID,
		ProviderTypeID: providerType.ID,
		Name:           dto.Name,
		Description:    *dto.Description,
		IsActive:       true,
		IsDefault:      dto.IsDefault,
		Config:         dto.Config,
		Credentials:    encryptedCreds,
		CreatedByID:    dto.CreatedByUserID,
	}

	// 4. Guardar en BD
	if err := uc.providerRepo.Create(ctx, provider); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to create provider")
		return nil, err
	}

	uc.log.Info(ctx).Uint("provider_id", provider.ID).Msg("Provider created successfully")
	return provider, nil
}

// UpdateProvider actualiza un proveedor existente
func (uc *useCase) UpdateProvider(ctx context.Context, id uint, dto *dtos.UpdateProviderDTO) error {
	uc.log.Info(ctx).Uint("provider_id", id).Msg("Updating invoicing provider")

	// Obtener proveedor existente
	provider, err := uc.providerRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrProviderNotFound
	}

	// Actualizar solo los campos proporcionados
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
		provider.IsDefault = *dto.IsDefault
	}

	// Encriptar credenciales si fueron proporcionadas
	if dto.Credentials != nil {
		encrypted, err := uc.encryption.Encrypt(dto.Credentials)
		if err != nil {
			uc.log.Error(ctx).Err(err).Msg("Failed to encrypt credentials")
			return err
		}
		provider.Credentials = encrypted
	}

	// Guardar cambios
	if err := uc.providerRepo.Update(ctx, provider); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to update provider")
		return err
	}

	uc.log.Info(ctx).Uint("provider_id", provider.ID).Msg("Provider updated successfully")
	return nil
}

// GetProvider obtiene un proveedor por ID
func (uc *useCase) GetProvider(ctx context.Context, id uint) (*entities.InvoicingProvider, error) {
	return uc.providerRepo.GetByID(ctx, id)
}

// ListProviders lista proveedores de un negocio
func (uc *useCase) ListProviders(ctx context.Context, businessID uint) ([]*entities.InvoicingProvider, error) {
	return uc.providerRepo.List(ctx, businessID)
}

// TestProviderConnection prueba la conexión con un proveedor
func (uc *useCase) TestProviderConnection(ctx context.Context, id uint) error {
	// Obtener proveedor
	provider, err := uc.providerRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrProviderNotFound
	}

	// Desencriptar credenciales
	credentials, err := uc.encryption.Decrypt(provider.Credentials)
	if err != nil {
		return errors.ErrDecryptionFailed
	}

	// Validar con el cliente del proveedor
	return uc.providerClient.ValidateCredentials(ctx, credentials)
}
