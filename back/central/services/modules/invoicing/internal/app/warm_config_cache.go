package app

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/log"
)

// WarmConfigCache pre-carga todas las configuraciones de facturación activas en Redis al iniciar el servidor.
// Cachea bajo CADA integration_id que tenga asociado, y también bajo el business_id del negocio.
func (uc *useCase) WarmConfigCache(ctx context.Context) error {
	ctx = log.WithFunctionCtx(ctx, "WarmConfigCache")

	uc.log.Info(ctx).Msg("🔥 Starting config cache warming for invoicing...")

	// 1. Obtener todas las configuraciones activas de BD (con ConfigIntegrations preloaded)
	configs, err := uc.repo.ListAllActiveConfigs(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Failed to load active configs for cache warming")
		return err
	}

	if len(configs) == 0 {
		uc.log.Info(ctx).Msg("⚠️ No active invoicing configs found - cache warming skipped")
		return nil
	}

	uc.log.Info(ctx).
		Int("total", len(configs)).
		Msg("📋 Active invoicing configs loaded from database")

	// 2. Cachear cada configuración por integration_id
	// GetConfigByIntegration hace la consulta completa (con preloads) y setea el caché
	successCount := 0
	errorCount := 0

	for _, config := range configs {
		for _, integrationID := range config.IntegrationIDs {
			cachedConfig, err := uc.repo.GetConfigByIntegration(ctx, integrationID)
			if err != nil {
				uc.log.Warn(ctx).
					Err(err).
					Uint("config_id", config.ID).
					Uint("integration_id", integrationID).
					Msg("⚠️ Failed to cache config for integration")
				errorCount++
				continue
			}

			if cachedConfig != nil {
				successCount++
			}
		}
	}

	// 3. Cachear por business_id usando los datos ya cargados de BD.
	// Se usa la primera config activa por negocio (ListAllActiveConfigs está ordenado por created_at DESC).
	seenBusinessIDs := make(map[uint]bool)
	businessSuccessCount := 0
	businessErrorCount := 0

	for _, config := range configs {
		if seenBusinessIDs[config.BusinessID] {
			continue // Solo cachear la primera config activa por negocio
		}
		seenBusinessIDs[config.BusinessID] = true

		// GetEnabledConfigByBusiness populará el caché vía read-through si no está ya
		if _, err := uc.repo.GetEnabledConfigByBusiness(ctx, config.BusinessID); err != nil {
			uc.log.Warn(ctx).
				Err(err).
				Uint("business_id", config.BusinessID).
				Uint("config_id", config.ID).
				Msg("⚠️ Failed to warm business config cache")
			businessErrorCount++
		} else {
			businessSuccessCount++
		}
	}

	uc.log.Info(ctx).
		Int("integration_success", successCount).
		Int("integration_errors", errorCount).
		Int("business_success", businessSuccessCount).
		Int("business_errors", businessErrorCount).
		Msg("✅ Config cache warming completed")

	return nil
}
