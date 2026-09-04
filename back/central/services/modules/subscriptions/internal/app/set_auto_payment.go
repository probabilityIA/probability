package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) SetAutoPaymentEnabled(ctx context.Context, businessID uint, enabled bool, actorUserID uint) error {
	if err := uc.repo.UpdateBusinessAutoPaymentEnabled(ctx, businessID, enabled); err != nil {
		return err
	}

	action := entities.AuditActionAutoPaymentDisabled
	description := "desactivo el pago automatico de la suscripcion"
	if enabled {
		action = entities.AuditActionAutoPaymentEnabled
		description = "activo el pago automatico de la suscripcion"
	}
	uc.recordAudit(ctx, businessID, actorUserID, action, description)
	return nil
}
