package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

// UpdateIntegration actualiza una integración existente
func (uc *IntegrationUseCase) UpdateIntegration(ctx context.Context, id uint, dto domain.UpdateIntegrationDTO) (*domain.Integration, error) {
	ctx = log.WithFunctionCtx(ctx, "UpdateIntegration")

	// Obtener integración existente
	existing, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationNotFound, err)
	}

	// Actualizar campos si se proporcionan
	if dto.Name != nil {
		existing.Name = *dto.Name
	}
	if dto.Code != nil {
		// Validar que el nuevo código no exista
		exists, err := uc.repo.ExistsIntegrationByCode(ctx, *dto.Code, existing.BusinessID)
		if err != nil {
			return nil, fmt.Errorf("error al verificar código: %w", err)
		}
		if exists && *dto.Code != existing.Code {
			return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationCodeExists, *dto.Code)
		}
		existing.Code = *dto.Code
	}
	if dto.StoreID != nil {
		existing.StoreID = *dto.StoreID
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	if dto.IsDefault != nil {
		existing.IsDefault = *dto.IsDefault
		if *dto.IsDefault {
			// Si se marca como default, desmarcar las demás
			if err := uc.repo.SetIntegrationAsDefault(ctx, id); err != nil {
				return nil, fmt.Errorf("error al marcar como default: %w", err)
			}
		}
	}
	if dto.Config != nil {
		existing.Config = *dto.Config
	}
	if dto.Credentials != nil {
		// DEBUG: Log credential keys being updated
		credKeys := make([]string, 0, len(*dto.Credentials))
		for k := range *dto.Credentials {
			credKeys = append(credKeys, k)
		}
		uc.log.Debug(ctx).
			Strs("credential_keys", credKeys).
			Uint("id", id).
			Msg("Credentials keys being updated in integration")

		// Serializar credenciales (se encriptarán en el repository)
		credentialsBytes, err := json.Marshal(*dto.Credentials)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsSerialize, err)
		}
		existing.Credentials = datatypes.JSON(credentialsBytes)
	}
	if dto.Description != nil {
		existing.Description = *dto.Description
	}

	// Solo actualizar UpdatedByID si es un ID válido (mayor que 0)
	if dto.UpdatedByID > 0 {
		existing.UpdatedByID = &dto.UpdatedByID
	}
	// Si UpdatedByID es 0, no actualizamos el campo (mantiene el valor existente o NULL)

	// ✅ NUEVO - Invalidar cache antes de actualizar
	if err := uc.cache.InvalidateIntegration(ctx, id); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to invalidate cache")
	}

	// Guardar cambios
	if err := uc.repo.UpdateIntegration(ctx, id, existing); err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al actualizar integración")
		return nil, fmt.Errorf("error al actualizar integración: %w", err)
	}

	// ✅ NUEVO - Re-cachear metadata actualizada
	integrationType, _ := uc.repo.GetIntegrationTypeByID(ctx, existing.IntegrationTypeID)

	configMap := make(map[string]interface{})
	if len(existing.Config) > 0 {
		json.Unmarshal(existing.Config, &configMap)
	}

	integrationTypeCode := ""
	if integrationType != nil {
		integrationTypeCode = integrationType.Code
	}

	cachedMeta := &domain.CachedIntegration{
		ID:                  existing.ID,
		Name:                existing.Name,
		Code:                existing.Code,
		Category:            existing.Category,
		IntegrationTypeID:   existing.IntegrationTypeID,
		IntegrationTypeCode: integrationTypeCode,
		BusinessID:          existing.BusinessID,
		StoreID:             existing.StoreID,
		IsActive:            existing.IsActive,
		IsDefault:           existing.IsDefault,
		Config:              configMap,
		Description:         existing.Description,
		CreatedAt:           existing.CreatedAt,
		UpdatedAt:           existing.UpdatedAt,
	}

	if err := uc.cache.SetIntegration(ctx, cachedMeta); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to cache updated metadata")
	}

	// ✅ NUEVO - Re-cachear credentials si cambiaron
	if dto.Credentials != nil {
		cachedCreds := &domain.CachedCredentials{
			IntegrationID: existing.ID,
			Credentials:   *dto.Credentials, // Ya están desencriptadas en el DTO
		}
		if err := uc.cache.SetCredentials(ctx, cachedCreds); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Failed to cache updated credentials")
		}
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Integración actualizada exitosamente")

	return existing, nil
}
