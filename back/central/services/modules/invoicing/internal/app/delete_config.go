package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// DeleteConfig elimina una configuración de facturación
func (uc *useCase) DeleteConfig(ctx context.Context, id uint) error {
	uc.log.Info(ctx).Uint("config_id", id).Msg("Deleting invoicing config")

	// Verificar que existe
	_, err := uc.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return errors.ErrConfigNotFound
	}

	// Eliminar
	if err := uc.repo.DeleteInvoicingConfig(ctx, id); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to delete config")
		return err
	}

	uc.log.Info(ctx).Uint("config_id", id).Msg("Config deleted successfully")
	return nil
}
