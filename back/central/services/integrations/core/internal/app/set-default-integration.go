package app

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/log"
)

// SetAsDefault marca una integración como default
func (uc *integrationUseCase) SetAsDefault(ctx context.Context, id uint) error {
	ctx = log.WithFunctionCtx(ctx, "SetAsDefault")

	if err := uc.repo.SetAsDefault(ctx, id); err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al marcar integración como default")
		return err
	}

	uc.log.Info(ctx).Uint("id", id).Msg("Integración marcada como default exitosamente")

	return nil
}
