package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// DeleteConfig elimina una configuración de facturación
func (uc *useCase) DeleteConfig(ctx context.Context, id uint) error {
	uc.log.Info(ctx).Uint("config_id", id).Msg("Deleting invoicing config")

	// Verificar que existe
<<<<<<< HEAD
	_, err := uc.configRepo.GetByID(ctx, id)
=======
	_, err := uc.repo.GetInvoiceByID(ctx, id)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if err != nil {
		return errors.ErrConfigNotFound
	}

	// Eliminar
<<<<<<< HEAD
	if err := uc.configRepo.Delete(ctx, id); err != nil {
=======
	if err := uc.repo.DeleteInvoicingConfig(ctx, id); err != nil {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
		uc.log.Error(ctx).Err(err).Msg("Failed to delete config")
		return err
	}

	uc.log.Info(ctx).Uint("config_id", id).Msg("Config deleted successfully")
	return nil
}
