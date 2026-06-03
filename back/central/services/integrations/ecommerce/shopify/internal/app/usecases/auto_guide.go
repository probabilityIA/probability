package usecases

import (
	"context"
	"fmt"
)

func (uc *SyncOrdersUseCase) SetAutoGenerateGuide(ctx context.Context, integrationID string, enabled bool) error {
	if _, err := uc.integrationService.GetIntegrationByID(ctx, integrationID); err != nil {
		return fmt.Errorf("error al obtener integracion: %w", err)
	}

	configUpdate := map[string]interface{}{
		"auto_generate_guide_enabled": enabled,
	}
	if err := uc.integrationService.UpdateIntegrationConfig(ctx, integrationID, configUpdate); err != nil {
		return fmt.Errorf("error al guardar la configuracion de auto-guia: %w", err)
	}
	return nil
}
