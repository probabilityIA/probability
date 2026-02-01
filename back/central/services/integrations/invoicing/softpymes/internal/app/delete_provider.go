package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/errors"
)

// DeleteProvider elimina un proveedor (soft delete)
func (uc *useCase) DeleteProvider(ctx context.Context, id uint) error {
	uc.log.Info(ctx).Uint("provider_id", id).Msg("Deleting Softpymes provider")

	// 1. Verificar que el proveedor existe
	provider, err := uc.providerRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrProviderNotFound
	}

	// 2. Verificar que no sea el proveedor por defecto
	if provider.IsDefault {
		return errors.ErrCannotDeleteDefault
	}

	// 3. Eliminar (soft delete en repository)
	if err := uc.providerRepo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to delete provider")
		return err
	}

	uc.log.Info(ctx).Uint("provider_id", id).Msg("Provider deleted successfully")
	return nil
}
