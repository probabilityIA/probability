package usecaseintegrations

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

// GetIntegrationByID obtiene una integración por su ID
func (uc *IntegrationUseCase) GetIntegrationByID(ctx context.Context, id uint) (*domain.Integration, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationByID")

	// ✅ NUEVO - Intentar cache primero
	cached, err := uc.cache.GetIntegration(ctx, id)
	if err == nil {
		uc.log.Debug(ctx).Uint("id", id).Msg("✅ Cache hit - metadata")

		// Convertir CachedIntegration a Integration
		configJSON, _ := json.Marshal(cached.Config)
		integration := &domain.Integration{
			ID:                id,
			Name:              cached.Name,
			Code:              cached.Code,
			Category:          cached.Category,
			IntegrationTypeID: cached.IntegrationTypeID,
			BusinessID:        cached.BusinessID,
			StoreID:           cached.StoreID,
			IsActive:          cached.IsActive,
			IsDefault:         cached.IsDefault,
			Config:            datatypes.JSON(configJSON),
			Description:       cached.Description,
			CreatedAt:         cached.CreatedAt,
			UpdatedAt:         cached.UpdatedAt,
		}
		// Cargar IntegrationType para mantener compatibilidad
		integrationType, _ := uc.repo.GetIntegrationTypeByID(ctx, cached.IntegrationTypeID)
		integration.IntegrationType = integrationType

		return integration, nil
	}

	// Cache miss - Leer de DB
	uc.log.Debug(ctx).Uint("id", id).Msg("⚠️ Cache miss - loading from DB")
	integration, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener integración")
		return nil, err
	}

	// ✅ NUEVO - Cachear para próxima vez
	configMap := make(map[string]interface{})
	if len(integration.Config) > 0 {
		json.Unmarshal(integration.Config, &configMap)
	}

	integrationTypeCode := ""
	if integration.IntegrationType != nil {
		integrationTypeCode = integration.IntegrationType.Code
	}

	cachedMeta := &domain.CachedIntegration{
		ID:                  integration.ID,
		Name:                integration.Name,
		Code:                integration.Code,
		Category:            integration.Category,
		IntegrationTypeID:   integration.IntegrationTypeID,
		IntegrationTypeCode: integrationTypeCode,
		BusinessID:          integration.BusinessID,
		StoreID:             integration.StoreID,
		IsActive:            integration.IsActive,
		IsDefault:           integration.IsDefault,
		Config:              configMap,
		Description:         integration.Description,
		CreatedAt:           integration.CreatedAt,
		UpdatedAt:           integration.UpdatedAt,
	}

	if err := uc.cache.SetIntegration(ctx, cachedMeta); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to cache metadata")
	}

	return integration, nil
}
