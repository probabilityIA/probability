package app

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/log"
)

// WarmConfigCache pre-carga todas las configuraciones de facturación activas en Redis al iniciar el servidor
// Esto evita el "cold start" en las primeras consultas de facturación automática
func (uc *useCase) WarmConfigCache(ctx context.Context) error {
	ctx = log.WithFunctionCtx(ctx, "WarmConfigCache")

	uc.log.Info(ctx).Msg("🔥 Starting config cache warming for invoicing...")

	// 1. Obtener todas las configuraciones activas de BD
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

	// 2. Cachear cada configuración
	successCount := 0
	errorCount := 0

	for _, config := range configs {
		// Usar GetConfigByIntegration para cachear (utiliza el mismo método que tiene cache-aside)
		// Esto asegura que el caché se pueble con el mismo formato que cuando se consulta normalmente
		cachedConfig, err := uc.repo.GetConfigByIntegration(ctx, config.IntegrationID)
		if err != nil {
			uc.log.Warn(ctx).
				Err(err).
				Uint("config_id", config.ID).
				Uint("integration_id", config.IntegrationID).
				Msg("⚠️ Failed to cache config")
			errorCount++
			continue
		}

		if cachedConfig != nil {
			successCount++
		}
	}

	uc.log.Info(ctx).
		Int("success", successCount).
		Int("errors", errorCount).
		Msg("✅ Config cache warming completed")

	return nil
}
