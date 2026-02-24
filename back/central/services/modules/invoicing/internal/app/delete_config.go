package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// DeleteConfig elimina permanentemente una configuración de facturación (hard delete)
func (uc *useCase) DeleteConfig(ctx context.Context, id uint) error {
	uc.log.Info(ctx).Uint("config_id", id).Msg("Deleting invoicing config")

	_, err := uc.repo.GetInvoicingConfigByID(ctx, id)
	if err != nil {
		return errors.ErrConfigNotFound
	}

	// Eliminar
	// Eliminar (el repositorio ya verifica existencia con Unscoped)
	if err := uc.repo.DeleteInvoicingConfig(ctx, id); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to delete config")
		return err
	}

	uc.log.Info(ctx).Uint("config_id", id).Msg("Config deleted successfully")
	return nil
}
