package app

import (
	"context"
)

// DeleteConfig elimina permanentemente una configuración de facturación (hard delete)
func (uc *useCase) DeleteConfig(ctx context.Context, id uint) error {
	uc.log.Info(ctx).Uint("config_id", id).Msg("Deleting invoicing config")

	// Eliminar (el repositorio ya verifica existencia con Unscoped)
	if err := uc.repo.DeleteInvoicingConfig(ctx, id); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to delete config")
		return err
	}

	uc.log.Info(ctx).Uint("config_id", id).Msg("Config deleted successfully")
	return nil
}
