package usecaseintegrations

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// WarmCache pre-carga todas las integraciones activas en Redis al iniciar el servidor
// Esto evita el "cold start" en las primeras consultas
func (uc *IntegrationUseCase) WarmCache(ctx context.Context) error {
	ctx = log.WithFunctionCtx(ctx, "WarmCache")

	uc.log.Info(ctx).Msg("🔥 Starting cache warming for integrations...")

	// 1. Obtener todas las integraciones activas de BD
	filters := domain.IntegrationFilters{
		IsActive: boolPtr(true),
		Page:     1,
		PageSize: 1000, // Cargar hasta 1000 integraciones activas
	}

	integrations, total, err := uc.repo.ListIntegrations(ctx, filters)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Failed to load integrations for cache warming")
		return err
	}

	if len(integrations) == 0 {
		uc.log.Info(ctx).Msg("⚠️ No active integrations found - cache warming skipped")
		return nil
	}

	uc.log.Info(ctx).
		Int64("total", total).
		Int("loaded", len(integrations)).
		Msg("📋 Integrations loaded from database")

	// 2. Cachear cada integración
	successCount := 0
	errorCount := 0

	for _, integration := range integrations {
		// 2.1 Cachear metadata
		if err := uc.cacheIntegrationMetadata(ctx, integration); err != nil {
			uc.log.Warn(ctx).
				Err(err).
				Uint("integration_id", integration.ID).
				Str("code", integration.Code).
				Msg("⚠️ Failed to cache metadata")
			errorCount++
			continue
		}

		// 2.2 Cachear credentials desencriptadas
		if len(integration.Credentials) > 0 {
			if err := uc.cacheIntegrationCredentials(ctx, integration); err != nil {
				uc.log.Warn(ctx).
					Err(err).
					Uint("integration_id", integration.ID).
					Str("code", integration.Code).
					Msg("⚠️ Failed to cache credentials")
				errorCount++
				continue
			}
		}

		successCount++
	}

	uc.log.Info(ctx).
		Int("success", successCount).
		Int("errors", errorCount).
		Msg("✅ Cache warming completed")

	return nil
}

// cacheIntegrationMetadata cachea metadata de una integración
func (uc *IntegrationUseCase) cacheIntegrationMetadata(ctx context.Context, integration *domain.Integration) error {
	configMap := make(map[string]interface{})
	if len(integration.Config) > 0 {
		if err := json.Unmarshal(integration.Config, &configMap); err != nil {
			return err
		}
	}

	integrationTypeCode := ""
	baseURL := ""
	baseURLTest := ""
	if integration.IntegrationType != nil {
		integrationTypeCode = integration.IntegrationType.Code
		baseURL = integration.IntegrationType.BaseURL
		baseURLTest = integration.IntegrationType.BaseURLTest
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
		IsTesting:           integration.IsTesting,
		Config:              configMap,
		Description:         integration.Description,
		CreatedAt:           integration.CreatedAt,
		UpdatedAt:           integration.UpdatedAt,
		BaseURL:             baseURL,
		BaseURLTest:         baseURLTest,
	}

	return uc.cache.SetIntegration(ctx, cachedMeta)
}

// cacheIntegrationCredentials cachea credentials desencriptadas
func (uc *IntegrationUseCase) cacheIntegrationCredentials(ctx context.Context, integration *domain.Integration) error {
	if len(integration.Credentials) == 0 {
		return nil
	}

	// Desencriptar credentials
	encryptedBytes, err := decodeEncryptedCredentials([]byte(integration.Credentials))
	if err != nil {
		return err
	}

	decrypted, err := uc.encryption.DecryptCredentials(ctx, encryptedBytes)
	if err != nil {
		return err
	}

	// Cachear credentials desencriptadas
	cachedCreds := &domain.CachedCredentials{
		IntegrationID: integration.ID,
		Credentials:   decrypted,
	}

	return uc.cache.SetCredentials(ctx, cachedCreds)
}

// boolPtr helper para crear puntero a bool
func boolPtr(b bool) *bool {
	return &b
}

// Nota: decodeEncryptedCredentials() ya existe en get-integration-by-type.go
// y es accesible desde cualquier archivo del mismo package
