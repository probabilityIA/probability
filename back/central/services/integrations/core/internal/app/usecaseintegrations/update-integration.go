package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

func (uc *IntegrationUseCase) UpdateIntegration(ctx context.Context, id uint, dto domain.UpdateIntegrationDTO) (*domain.Integration, error) {
	ctx = log.WithFunctionCtx(ctx, "UpdateIntegration")

	existing, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationNotFound, err)
	}

	oldCode := existing.Code
	oldStoreID := existing.StoreID
	if dto.Name != nil {
		existing.Name = *dto.Name
	}
	if dto.Code != nil {
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
		if *dto.StoreID != oldStoreID {
			if err := uc.ensureStoreIDNotInUse(ctx, *dto.StoreID, existing.IntegrationTypeID, existing.BusinessID, id); err != nil {
				return nil, err
			}
		}
		existing.StoreID = *dto.StoreID
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	if dto.IsDefault != nil {
		existing.IsDefault = *dto.IsDefault
		if *dto.IsDefault {
			if err := uc.repo.SetIntegrationAsDefault(ctx, id); err != nil {
				return nil, fmt.Errorf("error al marcar como default: %w", err)
			}
		}
	}
	if dto.Config != nil {
		existing.Config = *dto.Config
	}
	if dto.Credentials != nil {
		credKeys := make([]string, 0, len(*dto.Credentials))
		for k := range *dto.Credentials {
			credKeys = append(credKeys, k)
		}
		uc.log.Debug(ctx).
			Strs("credential_keys", credKeys).
			Uint("id", id).
			Msg("Credentials keys being updated in integration")

		credentialsBytes, err := json.Marshal(*dto.Credentials)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsSerialize, err)
		}
		existing.Credentials = datatypes.JSON(credentialsBytes)
	}
	if dto.Description != nil {
		existing.Description = *dto.Description
	}
	if dto.IsTesting != nil {
		existing.IsTesting = *dto.IsTesting
	}

	if dto.UpdatedByID > 0 {
		existing.UpdatedByID = &dto.UpdatedByID
	}

	if dto.Credentials == nil {
		existing.Credentials = nil
	}

	if dto.Credentials != nil {
		if err := uc.cache.InvalidateIntegration(ctx, id); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Failed to invalidate cache")
		}
	} else {
		if err := uc.cache.InvalidateMetadata(ctx, id); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Failed to invalidate metadata cache")
		}
	}

	if err := uc.repo.UpdateIntegration(ctx, id, existing); err != nil {
		if isUniqueStoreIDViolation(err) {
			uc.log.Warn(ctx).Uint("id", id).Msg("La base rechazo una cuenta de canal duplicada")
			return nil, domain.ErrIntegrationStoreIDInUse
		}
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al actualizar integración")
		return nil, fmt.Errorf("error al actualizar integración: %w", err)
	}

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
		IsTesting:           existing.IsTesting,
		Config:              configMap,
		Description:         existing.Description,
		CreatedAt:           existing.CreatedAt,
		UpdatedAt:           existing.UpdatedAt,
	}

	if oldCode != "" && oldCode != existing.Code {
		if err := uc.cache.InvalidateCodeIndex(ctx, oldCode); err != nil {
			uc.log.Warn(ctx).Err(err).Str("old_code", oldCode).Msg("Failed to invalidate old code index")
		}
	}

	if oldStoreID != "" && oldStoreID != existing.StoreID {
		if err := uc.cache.InvalidateStoreTypeIndex(ctx, oldStoreID, existing.IntegrationTypeID); err != nil {
			uc.log.Warn(ctx).Err(err).Str("old_store_id", oldStoreID).Msg("Failed to invalidate old store index")
		}
	}

	if err := uc.cache.SetIntegration(ctx, cachedMeta); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to cache updated metadata")
	}

	if dto.Credentials != nil {
		cachedCreds := &domain.CachedCredentials{
			IntegrationID: existing.ID,
			Credentials:   *dto.Credentials,
		}
		if err := uc.cache.SetCredentials(ctx, cachedCreds); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Failed to cache updated credentials")
		}
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Integración actualizada exitosamente")

	return existing, nil
}
