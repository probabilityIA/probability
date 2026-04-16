package usecaseintegrations

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/shared/log"
)

// UpdateLastSync actualiza el timestamp de última sincronización de una integración
func (uc *IntegrationUseCase) UpdateLastSync(ctx context.Context, integrationID string) error {
	ctx = log.WithFunctionCtx(ctx, "UpdateLastSync")

	// Convertir string a uint
	var id uint
	if _, err := fmt.Sscanf(integrationID, "%d", &id); err != nil {
		uc.log.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("Invalid integration ID format")
		return fmt.Errorf("invalid integration ID: %w", err)
	}

	// Actualizar en repositorio
	if err := uc.repo.UpdateLastSync(ctx, id, time.Now()); err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error updating last_sync")
		return err
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Last sync updated successfully")
	return nil
}
