package app

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/log"
)

// WarmConfigCache pre-carga todas las configuraciones de facturación activas en Redis al iniciar el servidor.
// Para cada config, cachea bajo CADA integration_id que tenga asociado.
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

	uc.log.Info(ctx).
		Int("success", successCount).
		Int("errors", errorCount).
		Msg("✅ Config cache warming completed")

	return nil
}
